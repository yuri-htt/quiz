/*
 CSVファイルを読み込んでクイズを出題し、
 回答者の正答率を記録する。
 クイズを読み込むファイルはproblems.csvをデフォルトとするが、フラグを使って任意に変更できる。
*/

package main

import (
  "encoding/csv"
  "flag"
  "fmt"
  "os"
)

func main() {
  // フラグを定義
  // - 引数1:オプション名
  // - 引数2:オプションのデフォルト値　
  // - 引数3:オプションの説明
  csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
  // コマンドラインを解析して定義されたフラグにセットする
  flag.Parse()
  // csvFilenameはファイル名である文字列へのポインタ
  file, err := os.Open(*csvFilename)
  if err != nil {
    exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
  }
  r := csv.NewReader(file)
  lines, err := r.ReadAll()
  if err != nil {
    exit("Failed to parse the provided CSV file")
  }
  problems := parseLines(lines)
  fmt.Println(problems)
}

// 引数: 文字列を要素とする構造体の配列
// 返却値:problem型の要素を格納する配列
func parseLines(lines [][]string) []problem {
    // 組み込みのmake関数を使って[]problem型で行数分だけ要素数を持つsliceを作成
    ret := make([]problem, len(lines))
    // forループでrangeを使ってsliceを１つずつ反復処理
    for i, line := range lines {
      ret[i] = problem{
        q: line[0],
        a: line[1],
      }
    }
    return ret
}

// 構造体problemを定義
type problem struct {
    q string
    a string
}

func exit(msg string) {
  fmt.Println(msg)
  os.Exit(1)
}
