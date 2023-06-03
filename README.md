# BlrecLiveMsgBot
基于Blrec录播姬和onebot v11的qq直播提醒机器人



## 如何部署





## 格式文本示例
* 开播  
  您关注的{user}正在直播：  
  {title}  
  [CQ:Image,{cover}]  
  https://live.bilibili.com/{roomid}
* 下播  
  您关注的{user}直播结束啦



## 监听webhook事件

| 事件                 | 类型                              | 备注     |
| -------------------- | --------------------------------- | -------- |
| 开播事件             | LiveBeganEvent                    | 正常处理 |
| 下播事件             | LiveEndedEvent                    | 正常处理 |
| 直播间信息改变事件   | RoomChangeEvent                   | 未处理   |
| 录制开始事件         | RecordingStartedEvent             | 未处理   |
| 录制完成事件         | RecordingFinishedEvent            | 未处理   |
| 录制取消事件         | RecordingCancelledEvent           | 未处理   |
| 视频文件创建事件     | VideoFileCreatedEvent             | 未处理   |
| 视频文件完成完成事件 | VideoFileCompletedEvent           | 未处理   |
| 弹幕文件创建事件     | DanmakuFileCreatedEvent           | 未处理   |
| 弹幕文件完成事件     | DanmakuFileCompletedEvent         | 未处理   |
| 原始弹幕文件创建事件 | RawDanmakuFileCreatedEvent        | 未处理   |
| 原始弹幕文件完成事件 | RawDanmakuFileCompletedEvent      | 未处理   |
| 视频后处理完成事件   | VideoPostprocessingCompletedEvent | 未处理   |
| 硬盘空间不足事件     | SpaceNoEnoughEvent                | 未处理   |
| 程序异常事件         | Error                             | 未处理   |



## 鸣谢

[Github - Bilibili Live Streaming Recorder (blrec) B站直播录播姬](https://github.com/acgnhiki/blrec)

[Github - ZeroBot 基于Onebot V11实现的go语言SDK](https://github.com/wdvxdr1123/ZeroBot)
