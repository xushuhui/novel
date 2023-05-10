# novel
你是一个go语言开发专家，帮我用go语言开发一个命令行应用。使用第三方包，减少代码量。举个例子，命令行输入novel help，输出所有命令参数和注释。输入novel search {test}，输出json内容 "{
"books": [
{
"_id": "343106",
"title": "从笑傲江湖到大明国师",
"author": "带刀校尉",
"shortIntro": "和左冷禅对过掌，和东方不败比过剑，卫央觉着江湖其实也没什么意思，也就下酒的大毛桃还可以，那就种一谷桃",
"cover": "60/a3/60a3740710ab4e87683c9c18c5ccaed0.jpg",
"cat": "修真仙侠",
"followerCount": "16",
"zt": "连载中",
"updated": "2022-05-07T16:06:51+00:00",
"lastchapter": "更新到:第1024章 鞑靼人可眼馋了"
}
]
}"。输入novel info {id}，输出json内容"{
"_id": "498806",
"author": "西湖遇雨",
"cover": "3d/d3/3dd340fcc0113e456fbe51e336539626.jpg",
"longIntro": "\r\n \r\n【狱中讲课，朱棣偷听后求我当国师】 ",
"title": "大明国师",
"zt": "连载中",
"cat": "都市·青春",
"wordCount": "1138942",
"retentionRatio": "55.50",
"followerCount": "0",
"updated": "2023-05-08T03:20:50+00:00",
"chaptersCount": "276",
"lastChapter": "五月求票！"
}"。输入novel chapter {id}，输出json内容"{
"_id": "498806",
"chaptercount": 276,
"chaptersUpdated": "2023-05-08T03:20:50+00:00",
"chapters": [
{
"title": "第1章 指点江山又不会改变什么",
"link": "498806/1",
"unreadble": false
}]}"。输入novel content {id} {chapter}，输出json内容"{
"title": "第199章 测算太阳？【求月票！】",
"body": "\n 第199章测算太阳？【求月票！】\n “日月为明.”\n 朱高煦今天的思维似乎异常活跃，他马上又问出了一个奇奇怪怪的问题。 "
}"。
输入novel download {id}，下载指定小说内容，先调用chapter方法获取所有chapters数量，然后for循环该
数量，调用content方法传入{id}和{chapter}，{chapter}就是每次循环的i+1，获取所有body，汇总一起后写入一个以小说标题命名的txt文件。使用"github.com/go-resty/resty/v2"
"github.com/urfave/cli/v2"，并且const baseURL = "https://api.aixdzs.com/"


