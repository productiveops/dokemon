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
  status: string
  state: string
  ports: IPort[]
}

export interface IImage {
  id: string
  name: string
  tag: string
  size: number
  created: number
}

export interface IVolume {
  driver: string
  name: string
}

export interface INetwork {
  id: string
  name: string
  driver: string
  scope: string
}

export interface IComposeLibraryItemHead {
  projectName: string
}

export interface IComposeLibraryItem {
  projectName: string
  definition: string
}

export interface INodeComposeItemHead {
  id: number
  projectName: string
  libraryProjectName: string
  status: string
}

export interface INodeComposeItem {
  id: number
  projectName: string
  libraryProjectName: string
  status: string
}

export interface INodeComposeContainer {
  name: string
  image: string
  service: string
  status: string
  ports: string
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
