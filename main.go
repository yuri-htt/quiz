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
  "strings"
  "time"
)

func main() {
  // フラグを定義
  // - 引数1:オプション名
  // - 引数2:オプションのデフォルト値　
  // - 引数3:オプションの説明
  csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
  timeLimit := flag.Int("limit", 5, "the time limit for the quiz in seconds")

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
  // Timerが終了すると現在の時刻がチャネルを使って構造体の要素Cに送信される(C <-chan Time)
  timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
  // チャネルからメッセージを受け取るまで待っている
  //<-timer.C

  correct := 0
problemloop:
  for i, p := range problems {
    fmt.Printf("Problem #%d: %s\n", i+1, p.q)
    answerCh := make(chan string)
    go func() {
      var answer string
      // Scanf:入力する変数の型を変換指定子に従って型変換する
      fmt.Scanf("%s\n", &answer)
      answerCh <- answer
    }()

    select {
    // 時間になったらメッセージを受け取り、その時点でスコアを出力
    case <-timer.C:
      // 時間になってチャネルがメッセージを受け取ったらfor文を抜ける
      break problemloop
    // 入力があればメッセージを受け取り、正解であればカウントアップ 
    case answer := <- answerCh:
      if answer == p.a {
        correct++
      }
    }
  }

  fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
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
        a: strings.TrimSpace(line[1]),
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
