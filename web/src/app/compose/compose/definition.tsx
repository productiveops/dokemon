import { useParams } from "react-router-dom"
import useNodeComposeItem from "@/hooks/useNodeComposeItem"
import ComposeDefinitionGitHub from "./definition-github"
import ComposeDefinitionLocal from "./definition-local"

export default function ComposeDefinition() {
  const { nodeId } = useParams()
  const { composeProjectId } = useParams()
  const { nodeComposeItem } = useNodeComposeItem(nodeId!, composeProjectId!)

  if (nodeComposeItem?.type === "github") {
    return <ComposeDefinitionGitHub />
  } else {
    return <ComposeDefinitionLocal />
  }
}
