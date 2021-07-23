# Copyright 2021 Linka Cloud  All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


FROM golang:alpine as builder

WORKDIR /go/go.linka.cloud/go-repo

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o go-repo .

FROM alpine

RUN apk add ca-certificates

COPY --from=builder /go/go.linka.cloud/go-repo/go-repo /usr/bin/

USER nobody

EXPOSE 8888

ENTRYPOINT ["/usr/bin/go-repo"]
