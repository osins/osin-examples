cd src 
go mod vendor
cd ..
docker build --no-cache -t uhub.service.ucloud.cn/humanrisk/sso .
docker run -it uhub.service.ucloud.cn/humanrisk/sso date
docker push uhub.service.ucloud.cn/humanrisk/sso
# docker save uhub.service.ucloud.cn/humanrisk/sso -o .\humanrisk.cn.sso.tar