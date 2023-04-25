#GOOS：目标平台的操作系统（darwin、freebsd、linux、windows）
#GOARCH：目标平台的体系架构（386、amd64、arm）
#交叉编译不支持 CGO 所以要禁用它

CGO_ENABLED=0 GOOS=windows GOARCH=arm go build