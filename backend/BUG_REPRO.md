# Bug 复现说明

## Bug 是什么

食谱推荐把过期食品误判为临期食品，统计报表把浪费金额和 TopWasted 统计到了临期食品上而漏掉过期食品，名称拼接会丢掉第一个名字。

## 如何触发

```bash
cd backend
go test ./internal/service/ -run TestFoodRecommendStats
```

## 错误信息

```
--- FAIL: TestFoodRecommendStats
    recommend foods=..., want only 临期牛奶
    waste amount=..., want 15 and 1
    JoinNames = ..., want 张三、李四、王五
```
