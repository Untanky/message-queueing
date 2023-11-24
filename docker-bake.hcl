variable "TAG" {
  default = "latest"
}

group "default" {
  targets = ["message-queueing"]
}

target "message-queueing" {
  dockerfile = "Dockerfile"
  platforms = ["linux/amd64", "linux/arm64"]
  tags = [
    "ghcr.io/untanky/message-queueing:${TAG}"
  ]
}
