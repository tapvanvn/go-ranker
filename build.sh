DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

arg_num=$#
if [ $arg_num -ne 1 ];then 
    echo "Invalid call"
    echo "syntax : ./build.sh <alpine | arm64v8>"
    exit 1
fi
arch=$1

tag=$(<./version.txt)-${arch}

server_url=tapvanvn

docker build -t $server_url/ranker:$tag  -f docker/$arch.dockerfile ./

docker push $server_url/ranker:$tag