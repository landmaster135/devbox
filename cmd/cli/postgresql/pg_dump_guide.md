# pg_dump ガイド

## pg_dump のインストール

```bash
sudo apt update

# if not installed
sudo apt install -y curl ca-certificates gnupg
```

PostgreSQL 17 のリポジトリを追加
```bash
curl -fsSL https://www.postgresql.org/media/keys/ACCC4CF8.asc \
  | sudo gpg --dearmor -o /usr/share/keyrings/postgresql.gpg

echo "deb [signed-by=/usr/share/keyrings/postgresql.gpg] https://apt.postgresql.org/pub/repos/apt $(. /etc/os-release && echo $VERSION_CODENAME)-pgdg main" \
  | sudo tee /etc/apt/sources.list.d/pgdg.list > /dev/null
```

目的のバージョンをインストール
```bash
sudo apt update

sudo apt install -y postgresql-client-17
```

実際に 17 を使うよう確認
```bash
type -a pg_dump

pg_dump --version
```

まだ 16 が優先される場合（PATHを17優先）
```bash
echo 'export PATH="/usr/lib/postgresql/17/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
pg_dump --version
```
