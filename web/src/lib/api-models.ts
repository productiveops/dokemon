// Request

export interface INodeCreateRequest {
  name: string
  environmentId: number | null
}

export interface INodeUpdateRequest {
  name: string
  environmentId: number | null
}

export interface INodeContainerBaseUrlUpdateRequest {
  containerBaseUrl: string
}

export interface ISettingUpdateRequest {
  value: string
}

// Response

export interface IPageResponse<T> {
  items: T[]
  pageNo: number
  pageSize: number
  totalRows: number
}

export interface ISetting {
  id: string
  value: string
}

export interface ICredentialHead {
  id: number
  name: string
  service?: string
  type: string
  userName?: string
}

export interface ICredential {
  id: number
  name: string
  service?: string
  type: string
  userName?: string
  secret: string
}

export interface IEnvironmentHead {
  id: number
  name: string
}

export interface IEnvironment {
  id: number
  name: string
}

export interface INodeHead {
  id: number
  name: string
  agentVersion: string
  environment: string
  online: boolean
  registered: boolean
  containerBaseUrl: string
}

export interface INode {
  id: number
  name: string
  environmentId: number
  containerBaseUrl: string
}

export interface IPort {
  ip: string
  privatePort: number
  publicPort: number
  type: string
}

export interface IContainer {
  id: string
  name: string
  image: string
  status: string
  state: string
  ports: IPort[]
  stale: string
}

export interface IImage {
  id: string
  name: string
  tag: string
  size: number
  dangling: boolean
  created: number
  inUse: boolean
}

export interface IVolume {
  driver: string
  name: string
  inUse: boolean
}

export interface INetwork {
  id: string
  name: string
  driver: string
  scope: string
  inUse: boolean
}

export interface IComposeLibraryItemHead {
  id?: number
  projectName: string
  type: string
}

export interface IFileSystemComposeLibraryItem {
  projectName: string
  definition: string
}

export interface IGitHubComposeLibraryItem {
  id: number
  credentialId: number
  projectName: string
  url: string
}

export interface INodeComposeItemHead {
  id: number
  projectName: string
  type: string
  libraryProjectId?: number
  libraryProjectName: string
  status: string
  stale: string
}

export interface INodeComposeItem {
  id: number
  projectName: string
  type: string
  libraryProjectId?: number
  libraryProjectName: string
  url?: string
  credentialId?: number
  definition?: string
  status: string
}

export interface INodeComposeVariable {
  id: number
  name: string
  isSecret: boolean
  value: string
}

export interface INodeComposeContainer {
  id: string
  name: string
  image: string
  service: string
  status: string
  state: string
  ports: string
  stale: string
}

export interface IComposeDefinition {
  projectName: string
  definition: string
}

export interface IComposeDefinition {
  projectName: string
  definition: string
}

export interface IVariable {
  id: number
  name: string
  isSecret: boolean
}

export interface IVariableHead {
  id: number
  name: string
  isSecret: boolean
  values: { [key: string]: string }
}
