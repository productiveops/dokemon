import { Terminal } from "@xterm/xterm"
import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export const CLASSES_TABLE_ACTION_ICON = "w-4 text-slate-900 dark:text-white"
export const CLASSES_CLICKABLE_TABLE_ROW =
  "cursor-pointer hover:bg-slate-50 dark:hover:bg-slate-800"

export function convertByteToMb(size: number) {
  return (size / (1000 * 1000)).toFixed(2) + " MB"
}

export function trimString(u: unknown) {
  return typeof u === "string" ? u.trim() : u
}

export function getContainerUrlFromPortMapping(
  portMapping: string,
  containerBaseUrl: string | null
) {
  const lastColonIndex = portMapping.lastIndexOf(":")
  const parts = portMapping.substring(lastColonIndex + 1).split("-")

  if (parts.length <= 1) {
    return null
  }

  let publicPort = parts[0]
  let hostname = portMapping.substring(0, lastColonIndex)

  let baseUrl = containerBaseUrl
  if (hostname === "0.0.0.0" || hostname == "::") {
    if (!baseUrl) {
      baseUrl = `${location.protocol}//${location.hostname}`
    }
    hostname = location.hostname
  } else {
    baseUrl = `${location.protocol}//${hostname}`
  }

  return `${baseUrl}:${publicPort}`
}

export function recreateTerminalElement(
  containerId: string,
  elementId: string
) {
  const elContainer = document.getElementById(containerId)
  if (elContainer) {
    const el = document.getElementById(elementId)
    if (el) {
      elContainer.removeChild(el)
    }
    const newEl = document.createElement("div")
    newEl.setAttribute("id", "terminal")
    elContainer.appendChild(newEl)
    return newEl
  }

  return null
}

export function newTerminal(convertToEol?: boolean) {
  if (convertToEol === undefined || convertToEol === null) {
    convertToEol = true
  }

  return new Terminal({
    theme: {
      background: "#0f172a",
    },
    fontFamily:
      'ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, "Liberation Mono", monospace',
    fontWeight: 100,
    cursorBlink: true,
    allowProposedApi: true,
    convertEol: convertToEol,
    rows: 28,
  })
}

export const REGEX_IDENTIFIER = /^[a-zA-Z0-9][a-zA-Z0-9_-]*$/
export const REGEX_IDENTIFIER_MESSAGE =
  "Only alphabets, digits, _ and -. Must start with an alphabet or digit."
