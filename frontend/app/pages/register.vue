<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { UserIcon, MailIcon, LockIcon } from 'lucide-vue-next'

definePageMeta({
  layout: 'auth'
})

const authStore = useAuthStore()
const router = useRouter()

const name = ref('')
const email = ref('')
const password = ref('')
const loading = ref(false)
const errorMessage = ref('')

async function onSubmit() {
  if (!name.value || !email.value || !password.value) {
    errorMessage.value = 'Please fill in all fields'
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    await authStore.register(name.value, email.value, password.value)
    // Auto-login after successful registration
    await authStore.login(email.value, password.value)
    router.push('/')
  } catch (error: any) {
    errorMessage.value = error.data?.message || 'Registration failed. Email might already be taken.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div>
    <form class="mt-8 space-y-6" @submit.prevent="onSubmit">
      <div class="space-y-4">
        <div class="space-y-2">
          <Label for="name">Full Name</Label>
          <div class="relative">
            <UserIcon class="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
            <Input
              id="name"
              v-model="name"
              type="text"
              placeholder="John Doe"
              class="pl-9 w-full"
              required
            />
          </div>
        </div>

        <div class="space-y-2">
          <Label for="email">Email address</Label>
          <div class="relative">
            <MailIcon class="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
            <Input
              id="email"
              v-model="email"
              type="email"
              placeholder="you@example.com"
              class="pl-9 w-full"
              required
            />
          </div>
        </div>

        <div class="space-y-2">
          <Label for="password">Password</Label>
          <div class="relative">
            <LockIcon class="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
            <Input
              id="password"
              v-model="password"
              type="password"
              placeholder="••••••••"
              class="pl-9 w-full"
              required
            />
          </div>
        </div>
      </div>

      <div v-if="errorMessage" class="text-sm text-red-500 bg-red-500/10 p-3 rounded-lg border border-red-500/20 text-center">
        {{ errorMessage }}
      </div>

      <div>
        <Button
          type="submit"
          class="w-full justify-center"
          size="lg"
          :disabled="loading"
        >
          <span v-if="loading" class="mr-2 animate-spin">⟳</span>
          Create account
        </Button>
      </div>
    </form>

    <div class="mt-6 text-center text-sm">
      <span class="text-muted-foreground">Already have an account? </span>
      <NuxtLink to="/login" class="font-medium text-primary hover:underline">
        Sign in here
      </NuxtLink>
    </div>
  </div>
</template>
