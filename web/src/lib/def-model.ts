export interface IDef {
  version: string
  services: Record<string, IService>
  volumes: Record<string, IVolume>
  networks: Record<string, INetwork>
}

export interface IService {
  image: string
  ports: IPortMapping[]
  volumes: string[]
  networks: string[]
  environment: Record<string, string>
  labels: Record<string, string>
  restart: string // no, always, on-failure, unless-stopped
}

export interface IPortMapping {
  target: number // Container port
  host_ip: string // Host IP
  published: string // Host Port
  protocol: string // tcp, udp
  mode: string // host, ingress
}

export interface IVolume {}

export interface INetwork {}
