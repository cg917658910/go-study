package main

import (
	"fmt"
	"regexp"
)

func extractCaseItems(text string) []string {
	// 正则表达式匹配 case '$item' : 模式
	re := regexp.MustCompile(`case\s+'([^']+)'\s*:`)

	// 查找所有匹配项
	matches := re.FindAllStringSubmatch(text, -1)

	// 提取匹配的item值
	var items []string
	for _, match := range matches {
		if len(match) > 1 {
			items = append(items, match[1])
		}
	}

	return items
}

func main() {
	text := `
		protected function pay_log_action($redis,$id){
        $model = new ShortMsgModel();
        $list = $model->where("handle_time",'=',0)->limit($this->limit)->select();
        if( $list ){
            foreach ( $list as $item ){
                $item['up_bank_coin'] = true;
                $up_time = $redis->get('update_bank_coin_' . $item['bank_number']);
                if( $up_time > time() ){
                    $item['up_bank_coin'] = false;
                }
                $item['ramk'] = preg_replace('/\[.*?\]/', '', $item['ramk']);
                try{
                    switch ($item['type_name']){
                        case 'SCB':
                        case 'SCB读取':
                            $this->handle_scb($item,$model);
                            break;
                        case 'SCB通知':
                            $this->handle_scb_notify($item,$model);
                            break;
                        case 'SCB流水':
                            $this->handle_scb_water($item,$model);
                            break;
                        case 'TM流水':
                            $this->handle_tm_water($item,$model);
                            break;
                        case 'KTB':
                            $this->handle_ktb($item,$model);
                            break;
                        case 'KTBLine':
                            $this->handle_ktb_line($item,$model);
                            break;
                        case 'KTB通知':
                            $this->handle_ktb_notice($item,$model);
                            break;
                        case 'KTB流水':
                            $this->handle_ktb_water($item,$model);
                            break;
                        case 'KBANK通知':
                            $this->handle_kkr_notify($item,$model);
                            break;
                        case 'KBANK读取':
                            $this->handle_kkr_read($item,$model);
                            break;
                        case 'KBANK流水':
                            $this->handle_kkr_water($item,$model);
                            break;
                        case 'BBL流水':
                            $this->handle_bbl_water($item,$model);
                            break;
                        case 'BBL':
                            $this->handle_bbl($item,$model);
                            break;
                        case 'BAAC':
                            $this->handle_baac($item,$model);
                            break;
                        case 'TTB':
                        case 'TTB读取':
                        case 'TTB通知':
                            $this->handle_ttb($item,$model);
                            break;
                        case 'KKRLSCLI':
                            $this->handle_kkr($item,$model);
                            break;
                        case 'BAY':
                            $this->handle_bay($item,$model);
                            break;
                        case 'TM流水ios':
                            $this->handle_tm_ios_water($item,$model);
                            break;
                        case 'SwooleTM':
                            $this->handle_tm_swool($item,$model);
                            break;
                        case 'GSB通知':
                        case 'GSB读取':
                        case 'GSBLine':
                            $this->handle_gsb($item,$model);
                            break;
                        case 'python-SCBGH':
                            $this->handle_scbgh($item,$model);
                            break;
                        case 'python-KTBGH':
                            $this->handle_ktbgh($item,$model);
                            break;
                        case 'python-Kbankgh':
                            $this->handle_kbankgh($item,$model);
                            break;
                        case 'python-Ttbgh':
                            $this->handle_ttbgh($item,$model);
                            break;
                        case 'python-BAYGH':
                            $this->handle_baygh($item,$model);
                            break;
                        case 'tm_protocol_water':
                            $this->tm_protocol_water($item,$model);
                            break;
                        case 'scb_protocol_water':
                            $this->scb_protocol_water($item,$model);
                            break;
                        case 'ttb_protocol_water':
                            $this->ttb_protocol_water($item,$model);
                            break;
                        default:
                            $model->where('pkbsm','=',$item['pkbsm'])->update([
                                'handle_time' => time(),
                                'handle_res' => '不是所需类型数据',
                            ]);
                            break;
                    }
                }catch (\Exception $e){
                    //修改
                    $model->where('pkbsm','=',$item['pkbsm'])->update([
                        'handle_time' => time(),
                        'handle_res' => "出现异常:".$e->getMessage().",行数：".$e->getLine(),
                    ]);
                }
            }
        }
    }

	`

	items := extractCaseItems(text)

	fmt.Println("提取结果：")
	for _, item := range items {
		fmt.Println(item)
	}
}
