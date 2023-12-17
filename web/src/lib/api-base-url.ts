export default function apiBaseUrl() {
  let url = window.location.origin
  return `${url}/api/v1`
}

export function wsApiBaseUrl() {
  let host = window.location.host
  if (host == "localhost:5173") {
    host = "localhost:9090"
  }

  const protocol = window.location.protocol === "http:" ? "ws" : "wss"
  return `${protocol}://${host}/api/v1`
}
