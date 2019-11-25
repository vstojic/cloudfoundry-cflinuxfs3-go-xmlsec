IMAGE_TAG=vladicastojic/cloudfoundry-cflinuxfs3-go-xmlsec
#docker rmi $IMAGE_TAG -f
#docker pull $IMAGE_TAG
docker build -t $IMAGE_TAG -f Dockerfile.build-static .
# Push it into the repo
docker commit $(docker ps -laq) couldfoundry-cflinuxfs3-go-xmlsec
docker tag cloudfoundry-cflinuxfs3-go-xmlsec $IMAGE_TAG
docker push $IMAGE_TAG
docker run -d $IMAGE_TAG sh
CONTAINER_ID=$(docker ps -alq)
docker cp $CONTAINER_ID:/go/src/github.com/crewjam/go-xmlsec/xmldsig .
docker stop $CONTAINER_ID -t 0

