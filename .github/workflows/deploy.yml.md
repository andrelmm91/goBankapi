# name: Deploy to Prod

# on:
#   push:
#     branches: ["master"]

# jobs:
#   build:
#     name: build image 2
#     runs-on: ubuntu-latest

#     permissions:
#       id-token: write
#       contents: read

#     steps:
#       - name: check out code
#         uses: actions/checkout@v4

#       - name: Configure AWS credentials
#         uses: aws-actions/configure-aws-credentials@v4
#         with:
#           role-to-assume: >>to be insert the IAM deployment ROLE
#           aws-region: eu-central-1

#       - name: login to Amazon ECR
#         id: login-ecr
#         uses: aws-actions/amazon-ecr-login@v2

#       - name: Build, tag, and push docker image to Amazon ECR
#         env:
#           REGISTRY: ${{ steps.login-ecr.outputs.registry }}
#           REPOSITORY: simplebank
#           IMAGE_TAG: ${{ github.sha }}
#         run: |
#           docker build -t $REGISTRY/$REPOSITORY:$IMAGE_TAG -f simplebank.dockerfile .
#           docker push $REGISTRY/$REPOSITORY:$IMAGE_TAG
