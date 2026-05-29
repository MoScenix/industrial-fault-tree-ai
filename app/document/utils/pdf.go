package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/MoScenix/industrial-fault-tree-ai/app/document/biz/model"
	"github.com/MoScenix/industrial-fault-tree-ai/app/document/conf"
	"github.com/ledongthuc/pdf"
)

const (
	defaultChunkSize    = 1000
	defaultChunkOverlap = 200
)

var errLowQualityPDFText = errors.New("parsed pdf text quality is too low")

func ResolvePDFPath(pdfID, fileName string) (string, error) {
	docDir := filepath.Join(conf.GetConf().Document.PDFDir, pdfID)
	if fileName != "" {
		candidate := filepath.Join(docDir, filepath.Base(fileName))
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	entries, err := os.ReadDir(docDir)
	if err != nil {
		return "", err
	}

	candidates := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		candidates = append(candidates, filepath.Join(docDir, entry.Name()))
	}
	sort.Strings(candidates)
	if len(candidates) == 0 {
		return "", os.ErrNotExist
	}
	return candidates[0], nil
}

func ParsePDFFile(pdfID, fileName string) (string, error) {
	pdfPath, err := ResolvePDFPath(pdfID, fileName)
	if err != nil {
		return "", err
	}
	return ParsePDFPath(pdfPath)
}

func ParsePDFPath(pdfPath string) (string, error) {
	type parseCandidate struct {
		text string
		err  error
	}

	candidates := []parseCandidate{
		{text: extractTextWithPDFToText(pdfPath)},
	}

	goPDFText, goPDFErr := extractTextWithGoPDF(pdfPath)
	candidates = append(candidates, parseCandidate{text: goPDFText, err: goPDFErr})

	bestText := ""
	bestQuality := pdfTextQuality{score: -1 << 30}
	var firstErr error
	for _, candidate := range candidates {
		if candidate.err != nil {
			if firstErr == nil {
				firstErr = candidate.err
			}
			continue
		}
		normalized := normalizePDFText(candidate.text)
		if normalized == "" {
			continue
		}
		quality := measurePDFTextQuality(normalized)
		if quality.score > bestQuality.score {
			bestText = normalized
			bestQuality = quality
		}
	}

	if bestText == "" {
		if firstErr != nil {
			return "", firstErr
		}
		return "", errLowQualityPDFText
	}
	if !bestQuality.usable() {
		return "", fmt.Errorf("%w: total=%d control=%d replacement=%d extended_latin=%d",
			errLowQualityPDFText, bestQuality.total, bestQuality.control, bestQuality.replacement, bestQuality.extendedLatin)
	}

	return bestText, nil
}

func extractTextWithPDFToText(pdfPath string) string {
	cmd := exec.Command("pdftotext", "-enc", "UTF-8", "-nopgbrk", "-layout", pdfPath, "-")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(output)
}

func extractTextWithGoPDF(pdfPath string) (string, error) {
	f, r, err := pdf.Open(pdfPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	plainText, err := r.GetPlainText()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(plainText); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func normalizePDFText(content string) string {
	content = strings.ReplaceAll(content, "\u0000", "")
	content = strings.ReplaceAll(content, "\f", "\n")
	lines := strings.Split(content, "\n")
	normalized := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		normalized = append(normalized, strings.Join(strings.Fields(line), " "))
	}
	return strings.TrimSpace(strings.Join(normalized, "\n"))
}

type pdfTextQuality struct {
	total         int
	control       int
	replacement   int
	extendedLatin int
	score         int
}

func (q pdfTextQuality) usable() bool {
	if q.total < 20 {
		return false
	}
	if q.control > 0 {
		return false
	}
	if q.replacement > 0 {
		return false
	}
	return q.score > 0
}

func measurePDFTextQuality(content string) pdfTextQuality {
	var q pdfTextQuality
	for _, r := range []rune(content) {
		if unicode.IsSpace(r) {
			continue
		}
		q.total++
		switch {
		case unicode.IsControl(r):
			q.control++
			q.score -= 6
		case r == unicode.ReplacementChar:
			q.replacement++
			q.score -= 6
		default:
			q.score++
			if isExtendedLatin(r) {
				q.extendedLatin++
				q.score -= 2
			}
			if unicode.Is(unicode.Han, r) {
				q.score++
			}
		}
	}
	return q
}

func isExtendedLatin(r rune) bool {
	return r >= 0x00C0 && r <= 0x024F
}

func BuildChunks(documentID, content string) []model.DocumentChunk {
	if content == "" {
		return nil
	}

	runes := []rune(content)
	chunks := make([]model.DocumentChunk, 0, len(runes)/defaultChunkSize+1)
	start := 0
	order := int64(1)

	for start < len(runes) {
		end := start + defaultChunkSize
		if end > len(runes) {
			end = len(runes)
		}

		chunkText := strings.TrimSpace(string(runes[start:end]))
		if chunkText != "" {
			chunks = append(chunks, model.DocumentChunk{
				ChunkID:    fmt.Sprintf("%s_%03d", documentID, order),
				DocumentID: documentID,
				Text:       chunkText,
				Page:       0,
				Order:      order,
			})
			order++
		}

		if end == len(runes) {
			break
		}

		start = end - defaultChunkOverlap
		if start < 0 {
			start = 0
		}
	}

	return chunks
}
