package utils

/**
 *
 *
 * @category   category
 * @package    Common
 * @subpackage Documentation\API
 * @author     zhangrubing  <zhangrubing@mioji.com>
 * @license    GPL https://mioji.com
 * @link       https://mioji.com
 */

//标准商品MODE
const (
	SCP_CUSTOM_TOUR       = 10000  //定制游
	SCP_FLIGHT            = 10001  //飞机
	SCP_TRAIN             = 10002  //火车
	SCP_FERRY             = 10003  //轮渡
	SCP_ZUCHE             = 10004  //租车
	SCP_BAOCHE            = 10005  //包车
	SCP_BUS               = 10006  //大巴
	SCP_IMAGE             = 10007  //图片资源
	SCP_PACKAGE_TOUR      = 10008  //跟团游
	SCP_ITINERARY         = 10009  //团游/行程
	SCP_DRIVE             = 10010  //驾车
	SCP_FUN               = 10011  //玩乐
	SCP_HOTEL             = 10012  //住宿
	SCP_GROUP_MEAL        = 10013  //团餐
	SCP_GUIDE             = 10014  //导游
	SCP_ATTRACTION_TICKET = 16384  //景点门票
	SCP_SHOW              = 32768  //演出/赛事
	SCP_SPECIAL_EVENTS    = 65536  //特色活动
	SCP_LOCAL_GROUP       = 524288 //当地参团
)

//POI MODE
const (
	//POI mode定义
	PM_CITY                  = 1       //城市
	PM_ATTRACTION            = 2       //景点
	PM_HOTEL                 = 4       //酒店
	PM_RESTAURANT            = 8       //餐厅
	PM_AIRPORT               = 16      //机场
	PM_RAILWAY_STATION       = 32      //火车站
	PM_CAR_RENTAL_POINT      = 64      //租车点
	PM_BUS_STATION           = 128     //汽车站
	PM_SHOPPING              = 256     //购物
	PM_COUNTRY               = 512     //国家
	PM_PROVINCE              = 1024    //省/州
	PM_FREE_ACTIVITIES       = 2048    //自由活动
	PM_OTHER_STATION         = 4096    //其他站点
	PM_FERRY_TERMINAL        = 8192    //轮渡站点
	PM_STANDARD_ROUTE        = 131072  //标准路线
	PM_CHARTER_POINTS        = 262144  //包车点
	PM_HOTEL_BRAND           = 1048576 //酒店品牌
	PM_HOTEL_BUSINESS_CIRCLE = 2097152 //酒店商圈
	PM_CUSTOM_CHARTER_POINTS = 4194304 //自定义包车区域
)

//ptype
const (
	P_DEFAULT = -1 //未定义
	P_NOLIMIT = 0  //不限
	P_ADULT   = 1  //成人
	P_CHILD   = 2  //儿童
	P_INFANT  = 3  //婴儿
)
