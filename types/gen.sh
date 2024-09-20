#!/usr/bin/sh
here=$(dirname $0)

flatc --go --gen-onefile --go-namespace model -o "$here/../internal/model" "$here/model.fbs"

flatc --go --gen-onefile --go-namespace services -o "$here/../internal/services" "$here/svc_group.fbs"
flatc --go --gen-onefile --go-namespace services -o "$here/../internal/services" "$here/svc_member.fbs"
