declare namespace API {
  type BaseResponseBoolean = {
    code?: number
    data?: boolean
    message?: string
  }

  type BaseResponseBytes = {
    code?: number
    data?: Blob
    message?: string
  }

  type BaseResponsePromptVO = {
    code?: number
    data?: PromptVO
    message?: string
  }

  type BaseResponseGraphEditState = {
    code?: number
    data?: GraphEditState
    message?: string
  }

  type BaseResponseGraphSuggestion = {
    code?: number
    data?: GraphSuggestionVO
    message?: string
  }

  type BaseResponseGraphVO = {
    code?: number
    data?: GraphVO
    message?: string
  }

  type BaseResponseLoginUserVO = {
    code?: number
    data?: LoginUserVO
    message?: string
  }

  type BaseResponseLong = {
    code?: number
    data?: number
    message?: string
  }

  type BaseResponsePageGraphMessageVO = {
    code?: number
    data?: PageGraphMessageVO
    message?: string
  }

  type BaseResponsePageGraphVO = {
    code?: number
    data?: PageGraphVO
    message?: string
  }

  type BaseResponsePageGraphVersionVO = {
    code?: number
    data?: PageGraphVersionVO
    message?: string
  }

  type BaseResponsePageUserVO = {
    code?: number
    data?: PageUserVO
    message?: string
  }

  type BaseResponseSaveResult = {
    code?: number
    data?: SaveResult
    message?: string
  }

  type BaseResponseString = {
    code?: number
    data?: string
    message?: string
  }

  type BaseResponseUser = {
    code?: number
    data?: User
    message?: string
  }

  type BaseResponseUserVO = {
    code?: number
    data?: UserVO
    message?: string
  }

  type BaseResponseWorkingGraph = {
    code?: number
    data?: WorkingGraphVO
    message?: string
  }

  type CreateGraphVersionRequest = {
    graphId?: number
    versionName?: string
  }

  type chatToModifyGraphParams = {
    graphId: number
    message: string
  }

  type DeleteGraphVersionRequest = {
    graphId?: number
    version?: string
  }

  type DeleteRequest = {
    id?: number
  }

  type DiscardWorkingGraphRequest = {
    graphId?: number
    version?: string
  }

  type downloadGraphParams = {
    graphId: number
    version?: string
    isTmp?: boolean
  }

  type getCurrentSuggestionParams = {
    graphId: number
    version?: string
  }

  type getGraphVOByIdParams = {
    id: number
  }

  type getPromptParams = {
    mode: number
  }

  type getUserByIdParams = {
    id: number
  }

  type getUserVOByIdParams = {
    id: number
  }

  type getWorkingGraphParams = {
    graphId: number
    version?: string
  }

  type GraphAddRequest = {
    graphName?: string
    description?: string
    cover?: string
  }

  type GraphEditState = {
    graphId?: number
    tmpReady?: boolean
    basedOnVersion?: string
    message?: string
  }

  type GraphMessageVO = {
    id?: number
    graphId?: number
    userId?: number
    role?: string
    content?: string
    createTime?: string
    updateTime?: string
  }

  type GraphQueryRequest = {
    pageNum?: number
    pageSize?: number
    sortField?: string
    sortOrder?: string
    id?: number
    graphName?: string
    description?: string
    userId?: number
  }

  type GraphSuggestionVO = {
    graphId?: number
    version?: string
    content?: string
    updateTime?: string
  }

  type GraphUpdateRequest = {
    id?: number
    graphName?: string
    description?: string
    cover?: string
  }

  type GraphVO = {
    id?: number
    graphName?: string
    description?: string
    cover?: string
    userId?: number
    currentVersion?: string
    hasTmp?: boolean
    createTime?: string
    updateTime?: string
  }

  type GraphVersionVO = {
    version?: string
    versionName?: string
    isCurrent?: boolean
    createTime?: string
    updateTime?: string
  }

  type listGraphMessageParams = {
    graphId: number
    pageSize?: number
    lastCreateTime?: string
  }

  type listGraphVersionParams = {
    graphId: number
  }

  type LoginUserVO = {
    id?: number
    userAccount?: string
    userName?: string
    userAvatar?: string
    userProfile?: string
    userRole?: string
    createTime?: string
    updateTime?: string
  }

  type PageGraphMessageVO = {
    records?: GraphMessageVO[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
  }

  type PageGraphVO = {
    records?: GraphVO[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
  }

  type PageGraphVersionVO = {
    records?: GraphVersionVO[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
  }

  type PageUserVO = {
    records?: UserVO[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
  }

  type RenameGraphVersionRequest = {
    graphId?: number
    version?: string
    versionName?: string
  }

  type PromptVO = {
    mode?: number
    content?: string
    updatedAt?: string
  }

  type SaveGraphRequest = {
    graphId?: number
    fromVersion?: string
    toVersion?: string
    remark?: string
    useTmp?: boolean
    content?: string
  }

  type SaveResult = {
    fromVersion?: string
    toVersion?: string
    message?: string
  }

  type ServerSentEventString = {
    d?: string
    message?: string
  }

  type StartEditRequest = {
    graphId?: number
    version?: string
  }

  type User = {
    id?: number
    userAccount?: string
    userPassword?: string
    userName?: string
    userAvatar?: string
    userProfile?: string
    userRole?: string
    editTime?: string
    createTime?: string
    updateTime?: string
    isDelete?: number
  }

  type UserAddRequest = {
    userName?: string
    userAccount?: string
    userAvatar?: string
    userProfile?: string
    userRole?: string
  }

  type UserLoginRequest = {
    userAccount?: string
    userPassword?: string
  }

  type UserQueryRequest = {
    pageNum?: number
    pageSize?: number
    sortField?: string
    sortOrder?: string
    id?: number
    userName?: string
    userAccount?: string
    userProfile?: string
    userRole?: string
  }

  type UserRegisterRequest = {
    userAccount?: string
    userPassword?: string
    checkPassword?: string
  }

  type UserUpdateRequest = {
    id?: number
    userName?: string
    userAvatar?: string
    userProfile?: string
  }

  type UserVO = {
    id?: number
    userAccount?: string
    userName?: string
    userAvatar?: string
    userProfile?: string
    userRole?: string
    createTime?: string
    updateTime?: string
  }

  type WorkingGraphVO = {
    graphId?: number
    version?: string
    isTmp?: boolean
    content?: string
  }

  type UpdatePromptRequest = {
    mode?: number
    content?: string
  }
}
