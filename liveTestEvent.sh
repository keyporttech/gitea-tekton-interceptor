#!/bin/bash

HEADERS='-H "X-Gitea-Delivery: 76cfa806-1c16-448b-b50d-2fe0789cf397" -H "X-Gitea-Event: push" -H "X-Gitea-Signature: 38ac8e178939fa0502ab0f616625fd975b2e6cb15d2cb2a565977c0900a2f7ca" -H "Content-Type: application/json"'

PAYLOAD=$(cat testEvent.json)

bash -c "curl $HEADERS --data @testEvent.json localhost:8080;"
