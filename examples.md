# utapi-go Usage Examples

Below are practical examples for all main features of the utapi-go module.

---

## 1. Delete Files

```go
utApi, _ := utapi.NewUtApi()
keys := []string{"your-file-key-1", "your-file-key-2"}
resp, err := utApi.DeleteFiles(keys)
if err != nil {
    panic(err)
}
fmt.Println("Deleted count:", resp.DeletedCount)
```

---

## 2. List Files

```go
utApi, _ := utapi.NewUtApi()
listResp, err := utApi.ListFiles(100, 0)
if err != nil {
    panic(err)
}
for _, file := range listResp.Files {
    fmt.Printf("File: %s, Size: %d\n", file.Name, file.Size)
}
```

---

## 3. Rename Files

```go
utApi, _ := utapi.NewUtApi()
updates := []utapi.RenameFileUpdate{{
    FileKey: "your-file-key",
    NewName: "new-filename.txt",
}}
renameResp, err := utApi.RenameFiles(updates)
if err != nil {
    panic(err)
}
fmt.Println("Renamed count:", renameResp.RenamedCount)
```

---

## 4. Get Usage Info

```go
utApi, _ := utapi.NewUtApi()
usage, err := utApi.GetUsageInfo()
if err != nil {
    panic(err)
}
fmt.Printf("Total bytes used: %d\n", usage.TotalBytes)
```

---

## 5. Get Presigned URL for Private File

```go
utApi, _ := utapi.NewUtApi()
fileKey := "your-private-file-key"
url, err := utApi.GetPresignedUrl(fileKey, 3600) // expires in 1 hour
if err != nil {
    panic(err)
}
fmt.Println("Presigned URL:", url)
```

---

## 6. Get App Info

```go
utApi, _ := utapi.NewUtApi()
appInfo, err := utApi.GetAppInfo()
if err != nil {
    panic(err)
}
fmt.Println("App ID:", appInfo.AppId)
```

---

## 7. Get Presigned URL for File Upload

```go
utApi, _ := utapi.NewUtApi()
files := []utapi.UploadFileInfo{{
    Name: "test.txt",
    Size: 123,
    Type: "text/plain",
}}
uploadResp, err := utApi.GetPresignedUploadUrl(files, "public-read")
if err != nil {
    panic(err)
}
presigned := uploadResp.Data[0]
fmt.Println("Upload URL:", presigned.Url)
```

---

## 8. Upload File to Presigned URL

```go
// presigned is utapi.PresignedPostURLs from GetPresignedUploadUrl
err := utapi.UploadFileToPresignedUrl("/path/to/test.txt", presigned)
if err != nil {
    panic(err)
}
fmt.Println("File uploaded successfully!")
```

---

## 9. Polling Upload Status (optional)

```go
// After uploading, you can poll presigned.PollingUrl to check upload status
resp, err := http.Get(presigned.PollingUrl)
if err != nil {
    panic(err)
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
fmt.Println(string(body))
```

---

For more details, see the GoDoc or the source code.
