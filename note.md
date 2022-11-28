---

---

server
chanX

client ---->
subChanX
go
{
select
<-chanX
<-dis
}

---
Item reward
---

pool
  itemX y%
  itemX y%

rules
  maxX y
  defaultItem

app
  poolCURD opLog ->
    no-concurency - update mem

  getCmd ->
    no-concurency
  ->
    out opLog

  interval -> checkpoint for each wal log