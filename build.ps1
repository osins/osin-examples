cd src 
go mod vendor
cd ..
docker build --no-cache -t humanrisk.cn/sso .
docker run -it humanrisk.cn/sso date
docker save humanrisk.cn/sso -o .\humanrisk.cn.sso.tar