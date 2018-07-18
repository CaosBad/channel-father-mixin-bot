<template>
<div class="home">
    <section>
        <h1>Welcome {{fullName}}</h1>
        <br>
        <b-loading :is-full-page="true" :active.sync="isLoading" :can-cancel="false"></b-loading>
    </section>
</div>
</template>

<script>


export default {
    data() {
        return {
            isLoading: false
        }
    },
    computed: {
        fullName() {
            return this.$store.state.account.full_name
        }
    }
}
</script>

<style lang="scss" scoped>
.home {
    text-align: center
}
</style>
