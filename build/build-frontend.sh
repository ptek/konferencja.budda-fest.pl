#!/bin/bash

set -euo pipefail

SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

export ZOLACMD=$(command -v zola)

if ! $(test -f "$ZOLACMD")
then
    (
        cd "$SCRIPTPATH"
        case "$(uname -s)" in
            Linux*)
                wget -nv https://github.com/getzola/zola/releases/download/v0.17.2/zola-v0.17.2-x86_64-unknown-linux-gnu.tar.gz
                ;;
            Darwin*)
                wget -nv https://github.com/getzola/zola/releases/download/v0.17.2/zola-v0.17.2-x86_64-apple-darwin.tar.gz
                ;;
            *)
                echo "Unsupported OS"
                exit
                ;;
        esac
        tar xvfz zola-*.tar.gz
    )
    ZOLACMD="./build/zola"
fi

$ZOLACMD build