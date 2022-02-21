DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

mkdir -p "$DIR/#temp/config/"

go build -o "$DIR/#temp/ranker" main.go 

rs=$?
if [ $rs -eq 0 ]; then 
    echo "SUCCESS"
    cp config_local.jsonc     "$DIR/#temp/config/config.jsonc"
    
    app="$DIR/#temp/ranker"
    kill $(lsof -t -i:9000)
    PORT=9000 $app

else
    echo "FAIL"
fi