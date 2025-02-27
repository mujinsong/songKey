# songKey
图外键  
In the actual production environment, using foreign keys is not recommended because of various troubles such as the cascade effect storm.
Therefore, the aim of this project is to provide foreign keys at the application level, and use a graph database to display the relationships directly between fields of different tables. It also shows how to quickly and concisely navigate from a certain field of one table to a certain field of another table. This empowers both newbies who haven't sorted out the table relationships and veterans who are too lazy to think about the table relationships with the ability to perform CRUD (Create, Read, Update, Delete) operations.
简易版介绍：实际生产环境中并不推荐使用外键，因为级联风暴之类的等等麻烦。

所以，本项目旨在提供应用层面的外键，并通过图数据库提供不同表字段直接的关系展示，以及如何快速简洁的从一个表的某字段到另一个表的某字段，为未捋清表关系新人以及懒得动脑思考表关系的老人提供crud赋能

例如：
表关系如下：  
![414da5b13c858722b2ba8ea005d1553e](https://github.com/mujinsong/songKey/assets/44770623/b04cc911-eeef-4378-8a9c-f7f923e59aa0)


查询从某一表的某字段如何到另一个表的字段：  
![f35dd4fa950584eaf176514653c61a5d](https://github.com/mujinsong/songKey/assets/44770623/324c6c9a-b3e5-4aec-a467-97a8c542b0c5)
