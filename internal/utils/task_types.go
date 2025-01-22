package utils

type TaskType string

const (
  TaskUploadFile TaskType = "UploadFile"
  TaskDownloadFile TaskType = "DownloadFile"
  TaskIssueGNRE TaskType = "IssueGNRE"
)
