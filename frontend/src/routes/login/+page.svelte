<script lang="ts">
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
  import { Alert, AlertDescription } from '$lib/components/ui/alert';
  import { Eye, EyeOff, Shield } from 'lucide-svelte';

  let username = '';
  let password = '';
  let showPassword = false;
  let isLoading = false;
  let error = '';

  async function handleLogin() {
    if (!username || !password) {
      error = 'Please enter both username and password';
      return;
    }

    isLoading = true;
    error = '';

    const result = await auth.login(username, password);
    
    if (result.success) {
      goto('/dashboard');
    } else {
      error = result.error || 'Login failed';
    }
    
    isLoading = false;
  }

  function handleKeyPress(event: KeyboardEvent) {
    if (event.key === 'Enter') {
      handleLogin();
    }
  }
</script>

<svelte:head>
  <title>Login - Smart City Surveillance</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 p-4">
  <Card class="w-full max-w-md">
    <CardHeader class="text-center">
      <div class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-primary-100">
        <Shield class="h-6 w-6 text-primary-600" />
      </div>
      <CardTitle class="text-2xl font-bold">Smart City Surveillance</CardTitle>
      <CardDescription>
        Sign in to access the surveillance dashboard
      </CardDescription>
    </CardHeader>
    
    <CardContent class="space-y-4">
      {#if error}
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      {/if}

      <div class="space-y-2">
        <Label for="username">Username</Label>
        <Input
          id="username"
          type="text"
          placeholder="Enter your username"
          bind:value={username}
          onkeypress={handleKeyPress}
          disabled={isLoading}
        />
      </div>

      <div class="space-y-2">
        <Label for="password">Password</Label>
        <div class="relative">
          <Input
            id="password"
            type={showPassword ? 'text' : 'password'}
            placeholder="Enter your password"
            bind:value={password}
            onkeypress={handleKeyPress}
            disabled={isLoading}
          />
          <Button
            type="button"
            variant="ghost"
            size="sm"
            class="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
            onclick={() => showPassword = !showPassword}
            disabled={isLoading}
          >
            {#if showPassword}
              <EyeOff class="h-4 w-4" />
            {:else}
              <Eye class="h-4 w-4" />
            {/if}
          </Button>
        </div>
      </div>

      <Button
        class="w-full"
        onclick={handleLogin}
        disabled={isLoading}
      >
        {#if isLoading}
          <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
        {/if}
        {isLoading ? 'Signing in...' : 'Sign In'}
      </Button>

      <div class="text-center text-sm text-gray-600">
        <p class="mb-2">Demo Credentials:</p>
        <div class="space-y-1 text-xs">
          <p><strong>SCS Operator:</strong> operator1 / password</p>
          <p><strong>Security Guard:</strong> guard1 / password</p>
        </div>
      </div>
    </CardContent>
  </Card>
</div> 