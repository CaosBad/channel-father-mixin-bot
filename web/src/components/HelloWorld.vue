<template>
<div class="container">
    <p>
        到开发者网站创建你的机器人
        <br>
        提交机器人的 ClientId, SessionId, PrivateKey
    </p>
    <form class="section" @submit.prevent="onSubmit">
        <b-field label="ClientId">
            <b-input v-model="clientId" maxlength="36"></b-input>
        </b-field>
        <b-field label="SessionId">
            <b-input v-model="sessionId" maxlength="36"></b-input>
        </b-field>
        <b-field label="Private Key">
            <b-input v-model="privateKey" maxlength="1024" type="textarea"></b-input>
        </b-field>
        <button class="button is-primary" v-bind:class="{'is-loading': loading}" ref="submitBtn" slot="trigger" type="submit">Submit</button>
    </form>
    <b-loading :is-full-page="false" :active.sync="isLoading" :can-cancel="false"></b-loading>
</div>
</template>

<script>
import api from '../api'

export default {
    name: 'HelloWorld',
    props: {
        bot: Object
    },
    data() {
        return {
            clientId: this.bot.client_id,
            sessionId: this.bot.session_id,
            privateKey: this.bot.private_key,
            loading: false
        }
    },
    methods: {
        onSubmit() {
            if (this.name === '' || this.sessionId === '' || this.privateKey === '') {
                this.$toast.open({
                    message: '不能为空',
                    type: 'is-danger'
                })
                return
            }
            this.loading = true
            api.pushKeys(this.bot.bot_id, this.clientId, this.sessionId, this.privateKey)
                .then((response) => {
                    this.loading = false
                    if (response.data.data) {
                        this.$toast.open({
                            duration: 6000,
                            message: '提交成功，请联系 ID:26596 开启服务',
                            type: 'is-success'
                        })
                    } else {
                        this.$toast.open({
                            message: '数据格式错误',
                            type: 'is-danger'
                        })
                    }
                }, (error) => {
                    console.log(error)
                })
            // this.$emit('submitKeys', this.model)
        }
    },
    computed: {
        isLoading() {
            return this.bot == "undefined"
        },
    },
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss" scoped>
h1 {
    font-weight: bold
}
</style>
