#!/bin/bash

# 1부터 100까지 반복
for i in {1..100}
do
    echo "--- 실행 횟수: $i ---"

    # 여기에 실행하고 싶은 명령어를 넣으세요.
    # go test를 하신다면 보통 go test ./... 형식을 사용합니다.
    go test -run TestRandao_Deploy

    # 바로 직전 명령어의 종료 상태가 0(성공)이 아니면 즉시 중단
    if [ $? -ne 0 ]; then
        echo "❌ $i 번째 실행에서 실패했습니다!"
        exit 1
    fi
done

echo "✅ 100회 실행 모두 성공했습니다!"
