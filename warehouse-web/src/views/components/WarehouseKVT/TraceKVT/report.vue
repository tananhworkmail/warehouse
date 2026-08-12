<template>
  <div class="wrapper">
    <div class="container">
      <div class="tabs-header">
        <el-button type="primary" class="back-btn" @click="goBack">
          Quay về
        </el-button>

        <div class="title">PO Information</div>
      </div>

      <el-tabs v-model="activeName" type="border-card">
        <el-tab-pane label="Customer PO" name="customer">
          <table
            class="my-table"
            style="border: 2px solid; padding: 5px"
            width="1024"
            rules="all"
            cellpadding="5"
          >
            <tr>
              <td colspan="3" align="Center">
                <h3>
                  {{ data.Customer.Head.Cp }} Purchase Order<br />NO：{{
                    data.Customer.Head.KHDDBH1
                  }}
                </h3>
              </td>
            </tr>
            <tr>
              <td align="center">Seller</td>
              <td align="center">Buyer</td>
              <td align="center">Ship To</td>
            </tr>
            <tr>
              <td width="400">
                FLOOURISH THRIVE DEVELOPMENTS LIMITED TAIWAN BRANCH <br />
                5F,NO 13,LANE 370,SEC 4,YA TAN ROAD,DA YA DISTRICT, TAICHUNG
                CITY,TAIWAN,R.O.C.<br />Vietname
              </td>
              <td width="400">{{ data.Customer.Head.BUYER }}</td>
              <td width="300">{{ data.Customer.Head.ShipTo }}</td>
            </tr>
          </table>
          <br />
          <br />
          <table
            class="my-table"
            style="border: 2px solid; padding: 5px"
            width="1024"
            rules="all"
            cellpadding="5"
          >
            <tr>
              <td align="center">Customer PO#</td>
              <td align="center">XieXing</td>
              <td align="center">XieMing</td>
              <td align="center">SKU</td>
              <td align="center">Description</td>
              <td align="center">QTY</td>
            </tr>
            <tr>
              <td>{{ data.Customer.Head.Customer_PO }}</td>
              <td>{{ data.Customer.Head.XieXing }}</td>
              <td>{{ data.Customer.Head.XieMing }}</td>
              <td>{{ data.Customer.Head.ARTICLE }}</td>
              <td>{{ data.Customer.Head.YSSM }}</td>
              <td>{{ data.Customer.Head.Qty }}</td>
            </tr>
          </table>
          <br />
          <table
            class="my-table"
            style="border: 2px solid; padding: 5px"
            width="1024"
            rules="all"
            cellpadding="5"
          >
            <tr>
              <td align="center">Size</td>
              <td align="center">Size Qty</td>
              <td align="center">Case Qty</td>
              <td align="center">Unit Cost</td>
            </tr>
            <tr v-for="(item, index) in data.Customer.Body" :key="index">
              <td>{{ item.Size }}</td>
              <td>{{ item.Size_Qty }}</td>
              <td>{{ item.Case_Qty }}</td>
              <td>{{ item.Unit_Case }}</td>
            </tr>
          </table>
        </el-tab-pane>
        <el-tab-pane label="BOM" name="bom">
          <table class="my-table">
            <tr>
              <td>
                <table class="my-table">
                  <tr>
                    <td>
                      <table
                        class="my-table"
                        style="border: 2px solid; padding: 5px"
                        width="724"
                        rules="all"
                        cellpadding="5"
                      >
                        <tr>
                          <td>PROJECT:</td>
                          <td>{{ data.Head.ZLBH }}</td>
                          <td>MODEL:</td>
                          <td>{{ data.Head.XieMing }}</td>
                          <td>ETD:</td>
                          <td>{{ data.Head.ShipDate }}</td>
                        </tr>
                        <tr>
                          <td>PO:</td>
                          <td>{{ data.Head.KHPO }}</td>
                          <td>ARTICLE:</td>
                          <td>{{ data.Head.ARTICLE }}</td>
                          <td>STD:</td>
                          <td>{{ data.Head.BZCC }}</td>
                        </tr>
                        <tr>
                          <td>Cstomer:</td>
                          <td>{{ data.Head.Kfjc }}</td>
                          <td>COLOR:</td>
                          <td>{{ data.Head.XieXing }}</td>
                          <td>LAST:</td>
                          <td>{{ data.Head.XTMH }}</td>
                        </tr>
                        <tr>
                          <td>COUNTRY:</td>
                          <td>{{ data.Head.Ywsm }}</td>
                          <td>COLOR DESC:</td>
                          <td>{{ data.Head.YSSM }}</td>
                          <td>CUT DIE:</td>
                          <td>{{ data.Head.DAOMH }}</td>
                        </tr>
                      </table>
                    </td>
                  </tr>
                  <tr>
                    <td>
                      <table
                        class="my-table"
                        style="border: 2px solid; padding: 5px"
                        width="724"
                        rules="all"
                        cellpadding="5"
                      >
                        <tr>
                          <td>SEO</td>
                          <td v-for="(item, index) in data.Size" :key="index">
                            {{ index + 1 }}
                          </td>
                        </tr>
                        <tr>
                          <td>SIZE RUN</td>
                          <td v-for="(item, index) in data.Size" :key="index">
                            {{ item.XXCC }}
                          </td>
                        </tr>
                        <tr>
                          <td>LAST</td>
                          <td v-for="(item, index) in data.Size" :key="index">
                            {{ item.XXCC }}
                          </td>
                        </tr>
                        <tr>
                          <td>OUTSOLE</td>
                          <td v-for="(item, index) in data.Size" :key="index">
                            {{ item.XXCC }}
                          </td>
                        </tr>
                        <tr>
                          <td>SOCK</td>
                          <td v-for="(item, index) in data.Size" :key="index">
                            {{ item.XXCC }}
                          </td>
                        </tr>
                        <tr>
                          <td>ORDER QTY</td>
                          <td v-for="(item, index) in data.Size" :key="index">
                            {{ item.QTY }}
                          </td>
                        </tr>
                      </table>
                    </td>
                  </tr>
                </table>
              </td>
              <td width="300" align="center">
                <img
                  :src="data.Head.IMGName"
                  alt=""
                  style="width: 220px; height: 220px"
                />
              </td>
            </tr>
          </table>
          <table
            class="my-table"
            style="border: 2px solid; padding: 5px; width: 1024px"
            rules="all"
            cellpadding="5"
          >
            <tr>
              <td>ID</td>
              <td>Name</td>
              <td>NAT NO</td>
              <td>DESCRIPTION</td>
              <td>UNIT</td>
              <td>SUPPLIER</td>
              <td>USAGE</td>
              <td>TOT QTY</td>
              <td>AREA</td>
            </tr>
            <tr v-for="(item, index) in data.Body" :key="index">
              <td>{{ item.BWBH }}</td>
              <td>{{ item.BWMC }}</td>
              <td>{{ item.CLBH }}</td>
              <td>{{ item.CLMC }}</td>
              <td>{{ item.DWBH }}</td>
              <td>{{ item.ZSYWJC }}</td>
              <td>{{ item.CLSL }}</td>
              <td>{{ item.SumQty }}</td>
              <td>{{ item.CQDH }}</td>
            </tr>
          </table>
        </el-tab-pane>
        <el-tab-pane label="Material PO" name="material">
          <div v-for="(item, index) in data.Orders" :key="index">
            <table class="my-table">
              <tr>
                <td>
                  <table class="my-table">
                    <tr>
                      <td width="100">Supplier</td>
                      <td width="140">{{ item.Head.Kfqm }}</td>
                      <td rowspan="4" width="400" align="center">
                        <H2
                          >material Order TY TNHH TY XUAN<BR />PHIEU DAT MUA
                          HANG</H2
                        >
                      </td>
                      <td width="120"></td>
                      <td width="120"></td>
                    </tr>
                    <tr>
                      <td>TEL</td>
                      <td>{{ item.Head.Dh }}</td>
                      <td></td>
                      <td>DATE</td>
                      <td>{{ item.Head.Dg_date }}</td>
                    </tr>
                    <tr>
                      <td>FAX</td>
                      <td>{{ item.Head.Cz }}</td>
                      <td></td>
                      <td>NO</td>
                      <td>{{ item.Head.Cghzbh }}</td>
                    </tr>
                  </table>
                </td>
              </tr>
              <tr>
                <td>
                  <table
                    class="my-table"
                    style="border: 2px solid; padding: 5px"
                    rules="all"
                    cellpadding="5"
                  >
                    <tr v-for="(item1, index) in item.Detail" :key="index">
                      <td width="30">{{ index + 1 }}</td>
                      <td width="120">{{ item1.Dgdh }}</td>
                      <td width="120">{{ item1.Clbh }}</td>
                      <td width="500">{{ item1.Ywpm }}</td>
                      <td width="60">{{ item1.Dwbh }}</td>
                      <td width="80">{{ item1.Zl_qty }}</td>
                      <td width="110">{{ item1.Zd_date }}</td>
                    </tr>
                  </table>
                </td>
              </tr>
            </table>
            <br />
            <hr />
            <br />
          </div>
          <div v-for="(material, index) in data.Material" :key="index">
            <table class="my-table">
              <tr>
                <td width="100">Supplier</td>
                <td width="100">{{ material.Title.ZSYWJC }}</td>
                <td rowspan="4" width="500" align="center">
                  <H2>CONG TY TNHH TY XUAN<BR />PHIEU DAT MUA HANG</H2>
                </td>
                <td width="100">Brand</td>
                <td width="120">BuyNo</td>
              </tr>
              <tr>
                <td>MatSupp</td>
                <td>{{ material.Title.ZSYWJC }}</td>
                <td></td>
                <td>Company</td>
                <td>{{ material.Title.GSBH }}</td>
              </tr>
              <tr>
                <td>TEL</td>
                <td>{{ material.Title.DH }}</td>
                <td></td>
                <td>DATE</td>
                <td>{{ material.Title.CGDate }}</td>
              </tr>
              <tr>
                <td>FAX</td>
                <td>{{ material.Title.CZ }}</td>
                <td></td>
                <td>NO</td>
                <td>{{ material.Title.CGNO }}</td>
              </tr>
            </table>

            <table
              class="my-table"
              style="border: 2px solid; padding: 5px"
              rules="all"
              cellpadding="5"
            >
              <tr v-for="(item, index) in material.Rows" :key="index">
                <td width="30">{{ index + 1 }}</td>
                <td width="120">{{ item.CLBH }}</td>
                <td width="500">{{ item.Ywpm }}</td>
                <td width="60">{{ item.DWBH }}</td>
                <td width="100" align="right">{{ item.Qty }}</td>
                <td width="120">{{ item.YQDate }}</td>
                <td></td>
              </tr>
            </table>
            <br />
            <hr />
            <br />
          </div>
        </el-tab-pane>
        <el-tab-pane label="Production" name="production">
          <table class="my-table">
            <tr>
              <td>
                <table
                  class="my-table"
                  align="center"
                  style="border: 2px solid; padding: 5px; width: 900px"
                  rules="all"
                  cellpadding="5"
                >
                  <tr>
                    <td colspan="4" align="center">
                      <h3>Production Of List</h3>
                    </td>
                  </tr>
                  <tr>
                    <td>XieXing</td>
                    <td>XieMing</td>
                    <td align="center">SKU</td>
                    <td>Description</td>
                  </tr>
                  <tr>
                    <td>{{ data.Head.XieXing }}</td>
                    <td>{{ data.Head.XieMing }}</td>
                    <td>{{ data.Head.ARTICLE }}</td>
                    <td>{{ data.Head.YSSM }}</td>
                  </tr>
                </table>
              </td>
            </tr>
            <tr>
              <td>
                <table
                  class="my-table"
                  v-if="data.Product.Items"
                  style="border: 2px solid; padding: 5px; width: 900px"
                  rules="all"
                  cellpadding="5"
                >
                  <tr>
                    <td>INDATE</td>
                    <td>Size</td>
                    <td>Qty</td>
                  </tr>
                  <tr v-for="(item, index) in data.Product.Items" :key="index">
                    <td>{{ item.INDATE }}</td>
                    <td>{{ item.XXCC }}</td>
                    <td>{{ item.Qty }}</td>
                  </tr>
                </table>
              </td>
            </tr>
            <tr>
              <td>
                <table
                  class="my-table"
                  v-if="data.Product.Items2"
                  style="border: 2px solid; padding: 5px; width: 900px"
                  rules="all"
                  cellpadding="5"
                >
                  <tr>
                    <td>INDATE</td>
                    <td>CARTONBAR</td>
                    <td>Qty</td>
                    <td>SGW</td>
                    <td>RGW</td>
                  </tr>
                  <tr v-for="(item, index) in data.Product.Items2" :key="index">
                    <td>{{ item.INDATE }}</td>
                    <td>{{ item.CARTONBAR }}</td>
                    <td>{{ item.Qty }}</td>
                    <td>{{ item.Sgw }}</td>
                    <td>{{ item.Rgw }}</td>
                  </tr>
                </table>
              </td>
            </tr>
          </table>
        </el-tab-pane>
        <el-tab-pane label="Shipping" name="shipping">
          <table class="my-table">
            <tr>
              <td colspan="5" align="center">
                <h3>
                  FLOOURISH THRIVE DEVELOPMENTS LIMITED TAIWAN BRANCH
                  <br />5F,NO 13,LANE 370,SEC 4,YA TAN ROAD,DA YA
                  DISTRICT,TAICHUNG CITY,TAIWAN,R.O.C. <br />COMMERCIAL INVOICE
                </h3>
              </td>
            </tr>
            <tr>
              <td>NO:</td>
              <td>LYN{{ data.Invoice.Head.INV_NO }}</td>
              <td width="200"></td>
              <td>Date:</td>
              <td>{{ data.Invoice.Head.Date }}</td>
            </tr>
            <tr>
              <td>Goods of:</td>
              <td>{{ data.Invoice.Head.Goods }} PRS OF FOOTWEAR</td>
              <td width="200"></td>
              <td></td>
              <td></td>
            </tr>
            <tr>
              <td>Shipped by :</td>
              <td>{{ data.Invoice.Head.Shipped }}</td>
              <td width="200"></td>
              <td></td>
              <td></td>
            </tr>
            <tr>
              <td>From:</td>
              <td>{{ data.Invoice.Head.FROM_WHERE }}</td>
              <td width="200"></td>
              <td>To:</td>
              <td>{{ data.Invoice.Head.TO_WHERE }}</td>
            </tr>
            <tr>
              <td colspan="5">
                <table
                  class="my-table"
                  style="width: 100%; border: 2px solid; padding: 5px"
                  rules="all"
                  cellpadding="5"
                >
                  <tr>
                    <td>PACK NO</td>
                    <td>Style Name</td>
                    <td>PO#</td>
                    <td>Color</td>
                    <td>Article#</td>
                    <td>Qty(PAIRS)</td>
                  </tr>
                  <tr v-for="(item, index) in data.Invoice.Body" :key="index">
                    <td>{{ item.Body.RYNO }}</td>
                    <td>{{ item.Body.XieMing }}</td>
                    <td>{{ item.Body.KHDDBH1 }}</td>
                    <td>{{ item.Body.YSSM }}</td>
                    <td>{{ item.Body.ARTICLE }}</td>
                    <td align="right">{{ item.Body.Pairs }}</td>
                  </tr>
                </table>
              </td>
            </tr>
            <tr>
              <td colspan="5"><hr /></td>
            </tr>
          </table>
          <br /><br />
          <table class="my-table">
            <tr>
              <td colspan="5" align="center">
                <h3>
                  FLOOURISH THRIVE DEVELOPMENTS LIMITED TAIWAN BRANCH
                  <br />5F,NO 13,LANE 370,SEC 4,YA TAN ROAD,DA YA
                  DISTRICT,TAICHUNG CITY,TAIWAN,R.O.C.
                  <p>PACKING LIST</p>
                </h3>
              </td>
            </tr>
            <tr>
              <td>INV NO:</td>
              <td>LYN{{ data.Invoice.Head.INV_NO }}</td>
              <td width="200"></td>
              <td>Date:</td>
              <td>{{ data.Invoice.Head.Date }}</td>
            </tr>
            <tr>
              <td>Goods of:</td>
              <td>{{ data.Invoice.Head.Goods }} PRS OF FOOTWEAR</td>
              <td width="200"></td>
              <td></td>
              <td></td>
            </tr>
            <tr>
              <td>Shipped by :</td>
              <td>{{ data.Invoice.Head.Shipped }}</td>
              <td width="200"></td>
              <td></td>
              <td></td>
            </tr>
            <tr>
              <td>From:</td>
              <td>{{ data.Invoice.Head.FROM_WHERE }}</td>
              <td width="200"></td>
              <td>To:</td>
              <td>{{ data.Invoice.Head.TO_WHERE }}</td>
            </tr>
            <tr>
              <td colspan="5">
                <div v-for="(item, index) in data.Invoice.Body" :key="index">
                  <table
                    class="my-table"
                    style="width: 100%; border: 2px solid"
                    rules="all"
                    cellpadding="5"
                  >
                    <tr>
                      <td>NO</td>
                      <td>{{ index + 1 }}</td>
                      <td>RY</td>
                      <td>{{ item.Body.RYNO }}</td>
                      <td>Color</td>
                      <td>{{ item.Body.YSSM }}</td>
                    </tr>
                    <tr>
                      <td>Style Name</td>
                      <td>{{ item.Body.XieMing }}</td>
                      <td>PO#</td>
                      <td>{{ item.Body.KHDDBH1 }}</td>
                      <td>Article</td>
                      <td>{{ item.Body.ARTICLE }}</td>
                    </tr>
                    <tr>
                      <td colspan="6">
                        <table
                          class="my-table"
                          style="border: 2px solid; padding: 5px; width: 100%"
                          rules="all"
                          cellpadding="5"
                        >
                          <tr>
                            <td>Number</td>
                            <td>Size</td>
                            <td>CTN</td>
                            <td>Cartons</td>
                            <td>Pairs</td>
                            <td>NW KGS</td>
                            <td>TNW KGS</td>
                            <td>RGW KGS</td>
                            <td>TRGW KGS</td>
                            <td>CBM L</td>
                            <td>CBM W</td>
                            <td>CBM H</td>
                          </tr>
                          <tr v-for="item1 in item.PackList" :key="item1">
                            <td>{{ item1.Number }}</td>
                            <td>{{ item1.SIZ }}</td>
                            <td>{{ item1.QTY }}</td>
                            <td>{{ item1.CTS }}</td>
                            <td>{{ item1.PAIRS }}</td>
                            <td>{{ item1.NW }}</td>
                            <td>{{ item1.TNW }}</td>
                            <td>{{ item1.GW }}</td>
                            <td>{{ item1.TGW }}</td>
                            <td>{{ item1.L }}</td>
                            <td>{{ item1.W }}</td>
                            <td>{{ item1.H }}</td>
                          </tr>
                          <tr>
                            <td colspan="2"></td>
                            <td colspan="2">Sub Total：{{ item.Total.CTS }}</td>
                            <td>{{ item.Total.Pairs }}</td>
                            <td></td>
                            <td>{{ item.Total.NW }}</td>
                            <td></td>
                            <td>{{ item.Total.GW }}</td>
                            <td colspan="3">
                              Sub Total CBM：{{ item.Total.CBM }}
                            </td>
                          </tr>
                        </table>
                      </td>
                    </tr>
                  </table>
                  <br />
                  <hr />
                  <br />
                </div>
              </td>
            </tr>
          </table>
        </el-tab-pane>
      </el-tabs>
    </div>
    <footer class="footer">© 2025 YIHSHUO. All rights reserved.</footer>
  </div>
</template>

<script>
export default {
  name: "TraceOrderReport",
};
</script>

<script setup>
import { ref, reactive, onMounted } from "vue";
import { useRouter } from "vue-router";
//import VuePdfApp from 'vue3-pdf-app';
import "vue3-pdf-app/dist/icons/main.css";

const API_URL = import.meta.env.VITE_API_URL;
const router = useRouter();
const goBack = () => {
  router.push("/warehouse-khode")
}
const data = reactive({
  Head: {},
  Body: [],
  Size: [],
  Material: [],
  Orders: [],
  Invoice: {
    Head: {},
    Body: [
      {
        Body: {},
        PackList: [],
        Total: {},
      },
    ],
  },
  Customer: {
    Head: {},
    Body: [],
  },
  Product: {
    Items: [],
    Items2: [],
  },
});
const activeName = ref("customer");
const state = reactive({
  productionInformationList: [],
});

async function getallinfomation() {
  const orderNo = router.currentRoute.value.params.orderNo;

  try {
    const res = await fetch(
      `${API_URL}/warehouse/bom_information?orderNo=${encodeURIComponent(
        orderNo
      )}`
    );

    if (!res.ok) {
      throw new Error(`HTTP error! status: ${res.status}`);
    }

    const result = await res.json();

    const d = result.data;

    data.Head = d.Bom?.Head || {};
    data.Size = d.Bom?.Size || {};
    data.Body = d.Bom?.Body || {};
    data.Material = d.Material || {};
    data.Orders = d.Orders || {};
    data.Invoice.Head = d.Invoice?.Head || {};
    data.Invoice.Body = d.Invoice?.Body || {};
    data.Customer.Head = d.Customer?.Head || {};
    data.Customer.Body = d.Customer?.Body || {};
    data.Product = d.Product || {};
  } catch (err) {
    console.error("Error fetching BOM information:", err);
    ElMessage.error("Error fetching BOM information");
  }
}

onMounted(() => {
  getallinfomation();
});
</script>

<style lang="scss" scoped>
.wrapper {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.container {
  flex: 1;
  padding: 20px;
  background-color: white;
  box-sizing: border-box;
}

.footer {
  text-align: center;
  padding: 10px 0;
  background-color: #000000;
  font-size: 14px;
}
:deep(.my-table) {
  color: #000 !important;
  border-color: #000 !important;
  background-color: #fff;
}
.tabs-header {
  display: flex;
  justify-content: center;  
  align-items: center;
  position: relative;    
  padding: 12px 0;
  background: white;
  border: 1px solid #dcdfe6;
  border-bottom: none;
  border-radius: 4px 4px 0 0;
  font-weight: bold;
  font-size: 18px;
}

.back-btn {
  position: absolute;
  left: 10px;
  top: 50%;
  transform: translateY(-50%);
}

.title {
  text-align: center;
  color: black;
}
</style>
