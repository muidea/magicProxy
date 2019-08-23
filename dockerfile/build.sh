#!/bin/sh

function buildImage()
{
    echo "build docker images..."
    docker build . > log.txt
    if [ $? -eq 0 ]; then
        echo "docker build success."
    else
        echo "docker build failed."
        exit 1
    fi
}

function tagImage()
{
    echo "tag docker images..."
    docker tag $1 registry.supos.ai/hdp/magicproxy
    if [ $? -eq 0 ];then
        echo "docker tag success."
    else
        echo "docker tag failed."
        exit 1
    fi
}

function rmiImage()
{
    echo "remove old docker images..."
    docker rmi $1
    if [ $? -eq 0 ];then
        echo "docker remove old images success."
    else
        echo "docker remove old images failed."
        exit 1
    fi    
}

#########################################################################################################
gopath=$GOPATH
#gopath=/home/jenkins/go
cp $gopath/bin/magicProxy ./

buildImage

oldID=$(docker images|grep registry.supos.ai/hdp/magicproxy|awk '{print $3}')
imageID=$(tail -1 log.txt |awk '{print $3}')

tagImage ${imageID}

#if [ ${oldID} != ${imageID} ];then
#    rmiImage ${oldID}
#fi

#rm -rf log.txt
#rm -rf magicProxy
