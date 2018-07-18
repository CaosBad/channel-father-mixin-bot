<template>
<div class="bot">
    <div v-if="getBot.bot_id">
        <h1>您已创建 Channel {{getBot.bot_id}}</h1>
        <h1>到期时间: {{getBot.expire_at | moment("YYYY-MM-DD")}}</h1>
        <br>
        <HelloWorld v-bind:bot="getBot" v-show="isKeyVisible" />
    </div>
    <div v-else>
        <h1>需要付费才能创建 Channel</h1>
        <br>
        <button class="button is-primary" v-on:click="onClickPay" v-bind:class="{'is-loading': loading}">点击支付</button>
        <br>
    </div>
</div>
</template>

<script>
// @ is an alias to /src
import uuidv4 from 'uuid/v4'
import HelloWorld from '@/components/HelloWorld.vue'
import Mixin from '@/mixin'
import api from '@/api'

export default {
    name: 'bot',
    data() {
        return {
            loading: false,
            isKeyVisible: true
        }
    },
    computed: {
        getBot() {
            return this.$store.getters.getBot
        }
    },
    components: {
        HelloWorld
    },
    methods: {
        waitForPayment: function (trace_id, redirect_to) {
            api.verifyPayment(trace_id)
                .then((response) => {
                    if (response.data.error) {
                        setTimeout(() => {
                            this.waitForPayment(trace_id, redirect_to)
                        }, 3000);
                        return true
                    } else if (response.data.data) {
                        this.loading = false
                        this.isKeyVisible = true
                        if (redirect_to != undefined) {
                            redirect_to.close();
                        }
                    }
                    this.$store.commit('SaveBot', response.data.data)
                }, (error) => {
                    console.log(error)
                })
        },
        onClickPay: function () {
            const trace_id = uuidv4()
            var redirect_to;
            if (Mixin.environment() == undefined) {
                redirect_to = window.open("");
            }
            if (Mixin.environment() == undefined) {
                redirect_to.location = 'https://mixin.one/pay?recipient=' + process.env.VUE_APP_CLIENT_ID + '&asset='+ process.env.VUE_APP_ASSET_ID +'&amount='+process.env.VUE_APP_ASSET_AMOUNT+'&trace=' + trace_id;
            } else {
                window.location.replace('mixin://pay?recipient=' + process.env.VUE_APP_CLIENT_ID + '&asset='+process.env.VUE_APP_ASSET_ID+'&amount='+process.env.VUE_APP_ASSET_AMOUNT+'&trace=' + trace_id);
            }
            this.loading = true
            setTimeout(() => {
                this.waitForPayment(trace_id, redirect_to)
            }, 3000);
        }

    },
    created: function () {
        api.getBot()
            .then((response) => {
                if (response.data.data) {
                    this.$store.commit('SaveBot', response.data.data)
                }
            }, (error) => {})
    }
}
</script>

<style lang="scss" scoped>
.bot {
    text-align: center
}
</style>
