# Bug 复现说明

## Bug 是什么

食品列表按到期日升序排序和分页偏移计算错误，CSV 导入会跳过只有两列的合法行，看板统计会给临期/过期列表各多塞一个空项，名称拼接会丢掉第一个名字。

## 如何触发

```bash
cd backend
go test ./internal/service/ -run TestFoodListOrderPage
```

## 错误信息

```
--- FAIL: TestFoodListOrderPage
    list order first=[], want 1
    page first=..., want ...
    import count=0 created=0, want 2 and 2
    expiring items len=..., want 1
    JoinNames = ..., want 张三、李四、王五
```
