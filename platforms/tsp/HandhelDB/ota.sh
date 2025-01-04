#!/bin/bash

# Nome do repositório no formato "owner/repo"
REPO="Handheld-Database/handheld-database-gui"

# Arquivo de versão local
VERSION_FILE="version.txt"

# URL da API do GitHub para tags
API_URL="https://api.github.com/repos/$REPO/tags"

# Nome do arquivo ZIP esperado na release
ZIP_FILE_NAME="trimui-HandheldDatabase.zip"

# Função para obter a versão local
get_local_version() {
    if [[ -f "$VERSION_FILE" ]]; then
        cat "$VERSION_FILE" | tr -d '[:space:]'
    else
        echo "v0.0.0" # Versão padrão caso não exista o arquivo
    fi
}

# Função para obter a versão mais recente do GitHub
get_latest_version() {
    curl -s "$API_URL" | jq -r '.[0].name'
}

# Função para baixar a release mais recente
download_latest_release() {
    local version="$1"
    local url="https://github.com/$REPO/releases/download/$version/$ZIP_FILE_NAME"

    echo "Baixando a release mais recente de: $url"
    curl -L -o latest_release.zip "$url"
}

# Função para limpar arquivos locais, exceto o ZIP e o ota.sh
clean_local_files() {
    echo "Limpando arquivos locais (exceto o ZIP e ota.sh)..."
    find . -mindepth 1 ! -name "latest_release.zip" ! -name "ota.sh" -exec rm -rf {} +
}

# Função para extrair a nova versão
extract_new_version() {
    echo "Extraindo nova versão..."
    unzip -o latest_release.zip -d temp_extract

    # Move todos os arquivos extraídos para a raiz do projeto
    mv temp_extract/*/* . 2>/dev/null || true

    # Remove a pasta temporária
    rm -rf temp_extract

    echo "Removendo o arquivo ZIP..."
    rm -f latest_release.zip
}

# Fluxo principal
local_version=$(get_local_version)
latest_version=$(get_latest_version)

echo "Versão local: $local_version"
echo "Versão mais recente: $latest_version"

if [[ "$local_version" == "$latest_version" ]]; then
    echo "Você já está na versão mais recente."
else
    echo "Atualizando para a versão mais recente..."
    download_latest_release "$latest_version"
    clean_local_files
    extract_new_version
    echo "$latest_version" > "$VERSION_FILE"
    echo "Atualização concluída para a versão $latest_version."
fi
