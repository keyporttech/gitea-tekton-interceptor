# Copyright 2020 Keyporttech Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#     http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#!/bin/bash

HEADERS='-H "X-Gitea-Delivery: 76cfa806-1c16-448b-b50d-2fe0789cf397" -H "X-Gitea-Event: push" -H "X-Gitea-Signature: 38ac8e178939fa0502ab0f616625fd975b2e6cb15d2cb2a565977c0900a2f7ca" -H "Content-Type: application/json"'

PAYLOAD=$(cat testEvent.json)

bash -c "curl $HEADERS --data @testEvent.json localhost:8080;"
