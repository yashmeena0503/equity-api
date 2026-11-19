import { useCallback, useState, type KeyboardEvent } from "react";
import {
  ReactFlow,
  Background,
  Controls,
  MiniMap,
  useEdgesState,
  useNodesState,
  type Edge,
  type Node,
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";
import { API_URL } from "../config"
import './NodeGraph.css'

export default function NodeGraph() {
  const [nodes, setNodes, onNodesChange] = useNodesState<Node>([])
  const [edges, setEdges, onEdgesChange] = useEdgesState<Edge>([])
  const [userInput, setUserInput] = useState("")
  const [status, setStatus] = useState("")

  // add user node
  // checks whether the user exists through the api
  // the api returns { "result": true|false }
  const addUserNode = useCallback(async () => {
    const unid = userInput.trim()
    if (!unid) {
      setStatus("enter a user unid")
      return
    }

    setStatus("checking user...")

    try {
      const response = await fetch(`${API_URL}/api/webapp/user/${encodeURIComponent(unid)}`)
      const responseText = await response.text()
      console.log("API response", {
        status: response.status,
        contentType: response.headers.get("content-type"),
        body: responseText,
      })
      if (!response.ok) { throw new Error(`HTTP ${response.status}`) }

      let data: { result: boolean }
      try { data = JSON.parse(responseText) }
      catch { throw new Error(`API returned invalid JSON: ${responseText}`) }
      if (!data.result) {
        setStatus("user does not exist")
        return
      }

      const nodeId = `user-${unid}`
      setNodes((currentNodes) => {
        // don't add duplicate user nodes
        if (currentNodes.some((node) => node.id === nodeId)) {
          setStatus("user already on graph")
          return currentNodes
        }

        const nodeNumber = currentNodes.length

        const newNode: Node = {
          id: nodeId,
          // simple positioning for now
          // later to be replaced by automatic graph layout
          position: {
            x: 100 + (nodeNumber %5) * 220,
            y: 100 + Math.floor(nodeNumber / 5) * 140,
          },
          data: {
            label: (
              <div className="graph-node">
                <div className="graph-node-type">user</div>
                <div className="graph-node-value">{unid}</div>
              </div>
            ),
          },
          type: "default",
        }
        return [...currentNodes, newNode]
      })
      setUserInput("")
      setStatus("user added")
    } catch (error) {
      console.error("failed to check user: ", error)
      setStatus("failed to check user")
    }
  }, [userInput, setNodes])

  // allow pressing enter in the input to add the user
  const handleUserInputKeyDown = (
    event: KeyboardEvent<HTMLInputElement>,
  ) => {
    if (event.key === "Enter") {
      void addUserNode()
    }
  }

  return (
    <div className="graph-page">
      {/* Toolbar */}
      <div className="graph-toolbar">
        <input
          type="text"
          placeholder="user unid"
          value={userInput}
          onChange={(event) => setUserInput(event.target.value)}
          onKeyDown={handleUserInputKeyDown}
        />

        <button onClick={() => void addUserNode()}>
          Add User
        </button>

        {status && <span>{status}</span>}
      </div>

      {/* React Flow graph */}
      <div className="graph-container">
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          fitView
        >
          <Background />
          <Controls />
          <MiniMap />
        </ReactFlow>
      </div>
    </div>
  )
}
