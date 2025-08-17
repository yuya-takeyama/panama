# PANAMA 変更設計（go-fuzzyfinder 採用版）

## 要点まとめ

* UI：**内蔵ファジーファインダ = github.com/ktr0731/go-fuzzyfinder** を採用
  → 外部バイナリ不要／クロスプラットフォーム👍
* 非TTY（CIやパイプ）は **stdio フォールバック**（番号選択 or `--first`）
* `peco` は完全ドロップ、`fzf` 依存もなし（オプションにも出さない）

## CLI 仕様（更新）

**コマンド構成**（変更なし）

* `panama select` … 内蔵 fuzzyfinder で選択 → **単一パス**出力
* `panama list` … 候補列挙（`--json`/`--format=...`）
* `panama init` / `panama version`

**共通フラグ**（変更点のみ）

* `--query, -q <str>` … 初期クエリ（fuzzyfinder に渡す）
* `--first` … 複数候補でも**最上位を即決**（非対話用途）
* （削除）`--picker` … なくす
* 既存：`--format path|cd|json`, `--max-depth`, `--no-cache`, `--silent` は継続

**使用例**

```bash
# いつもどおりインタラクティブ
panama select -q api

# 非TTYやスクリプトで
cd "$(panama select --first)"
```

## 設定ファイル仕様（更新）

探索はこれまでどおり `.panama.{yaml,toml,json}` を**上方探索**。

**変更点**

```yaml
# 旧: picker: auto   # ←削除
ui: fuzzyfinder       # fuzzyfinder | stdio
```

* 後方互換：もし `picker:` が残ってても無視（警告ログのみ）

## 内蔵 UI 実装方針

**パッケージ**：`internal/ui/fuzzyfinder`

* 提供IF：

  ```go
  type Item struct {
    Label       string
    Description string
    Path        string
    Score       float64
  }
  func Select(items []Item, query string) (idx int, err error)
  ```
* 実体：`go-fuzzyfinder` の `Find` を薄くラップ

  * プロンプト：`workspaces > `
  * ラベル：`Label`
  * サブ：`Description`
  * プレビュー：`Path`（折返し）
  * 事前クエリ：`WithQuery(query)` 相当（ライブラリがなければ、最初に type するだけでもOK）
* **TTY 判定**で UI/stdio を自動切替

  * `isatty` 系（`golang.org/x/term`）
  * 非TTY or `--first` → ソート済み先頭を返す
  * `--format=json` のときは list 相当もOK

**サンプル（中核部）**

```go
idx, err := fuzzyfinder.Find(
  items,
  func(i int) string { return items[i].Label },
  fuzzyfinder.WithPromptString("workspaces > "),
  fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
    if i < 0 || i >= len(items) { return "" }
    return items[i].Path
  }),
)
```

## アーキ更新

```
internal/ui/fuzzyfinder     # go-fuzzyfinder UI + stdioフォールバック
# internal/picker は削除
```

## 依存ライブラリ

* 追加：`github.com/ktr0731/go-fuzzyfinder`
* 既存：`cobra`, `doublestar`, `yaml/toml`, ほか
* （暗黙依存で）`golang.org/x/term`, `github.com/mattn/go-runewidth` 等が入る可能性あり

## 非TTY/CI の挙動

* `select --first`：スコア順の先頭を即出力（失敗時は終了コード1）
* `select`（TTYなし・`--first`なし）：番号選択の**stdio**モード
* `list --json`：UIなしで候補列挙

## テスト戦略（更新）

* **UI**：TTY 必須のためスモーク（起動→すぐ終了）を最小限。ロジックは**UI層の外側**で網羅。
* **stdio**：e2e で番号入力パスを網羅。
* **非TTY判定**：テーブル駆動＋環境変数で強制モード切替できるように（例：`PANAMA_FORCE_STDIO=1`）

## 実装タスク（差し替え）

1. `internal/ui/fuzzyfinder` を新規作成（上記IF）
2. Cobra の `select` で TTY 判定 → UI or stdio へ
3. `--first` 実装（スコア済コレクションから即決）
4. 設定 `ui: fuzzyfinder|stdio` を解釈（ただし TTYなしは stdio 強制）
5. 旧 `picker` 設定キーは**非推奨警告**のみ
6. README/設計書の fzf/peco 記述を削除・更新

---

## ちょい実装スケルトン（抜粋）

```go
// cmd/panama/select.go
func runSelect(cmd *cobra.Command, args []string) error {
  root := resolveRoot(opts)
  cfg  := config.Load(opts.Config, root)

  cands := pipeline.CollectScoreSort(root, cfg, opts) // 既存パイプライン

  if opts.First || !isTTY(os.Stdout.Fd()) || cfg.UI == "stdio" {
    if len(cands) == 0 { return fmt.Errorf("no candidates") }
    return output.Print(cands[0].Path, opts.Format)
  }

  items := toItems(cands)
  idx, err := ui.Select(items, opts.Query)
  if err != nil { return err }
  return output.Print(items[idx].Path, opts.Format)
}
```
