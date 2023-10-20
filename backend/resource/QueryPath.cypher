match (n) where id(n) in %v with collect(n) as st
match (n) where id(n) in %v with collect(n) as ed,st
match path=(qd)-[*%v]->(zd)
where qd in st and zd in ed
with path,relationships(path) as pts1 ,nodes(path) as nds1,st,ed
unwind nds1 as ndTemp
unwind pts1 as ptTemp
unwind path as phTemp
with apoc.coll.toSet(collect(ptTemp)) as pts,apoc.coll.toSet(collect(ndTemp)) as nds,st,ed,collect(phTemp) as ph
with [n in nds where size([pt in pts where endNode(pt)=n])>=n.degree]+st as filterNd,nds,ph
with [p in ph where ALL(n in nodes(p) where n in filterNd)] as ans
unwind ans as ansTemp
with nodes(ansTemp) as nds,relationships(ansTemp) as pts
return distinct(nds),pts