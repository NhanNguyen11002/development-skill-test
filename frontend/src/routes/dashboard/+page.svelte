<script lang="ts">
  import { onMount } from 'svelte';
  import { auth } from '$lib/stores/auth';
  import { alerts } from '$lib/stores/alerts';
  import { api } from '$lib/api';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Tabs, TabsContent, TabsList, TabsTrigger } from '$lib/components/ui/tabs';
  import { Avatar, AvatarFallback, AvatarImage } from '$lib/components/ui/avatar';
  import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from '$lib/components/ui/dropdown-menu';
  import { 
    Camera, 
    AlertTriangle, 
    Users, 
    Settings, 
    LogOut, 
    Bell,
    Video,
    MapPin,
    Clock
  } from 'lucide-svelte';
  import type {  Alert, Camera as CameraType } from '$lib/types';

  let cameraGrid = $state<CameraType[]>([])
  let isLoading = $state(true);
  let error = $state('');

  $inspect('cameraGrid', cameraGrid)
  onMount(async () => {
    await loadDashboard();
  });

  async function loadDashboard() {
    isLoading = true;
    error = '';

    try {
      const [gridResponse] = await Promise.all([
        api.getCameras()
      ]);
      console.log('response',gridResponse )
      if (gridResponse.data) {
        cameraGrid = gridResponse.data;
      } else {
        error = gridResponse.error || 'Failed to load camera grid';
      }
    } catch (err) {
      error = 'Failed to load dashboard data';
    } finally {
      isLoading = false;
    }
  }

  function handleLogout() {
    auth.logout();
  }

  function getSeverityColor(severity: Alert['severity']) {
    switch (severity) {
      case 'critical': return 'bg-red-100 text-red-800';
      case 'high': return 'bg-orange-100 text-orange-800';
      case 'medium': return 'bg-yellow-100 text-yellow-800';
      case 'low': return 'bg-green-100 text-green-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  }

  function getStatusColor(status: Alert['status']) {
    switch (status) {
      case 'pending': return 'bg-yellow-100 text-yellow-800';
      case 'acknowledged': return 'bg-blue-100 text-blue-800';
      case 'assigned': return 'bg-purple-100 text-purple-800';
      case 'resolved': return 'bg-green-100 text-green-800';
      case 'closed': return 'bg-gray-100 text-gray-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  }
</script>

<svelte:head>
  <title>Dashboard - Smart City Surveillance</title>
</svelte:head>

<div class="min-h-screen bg-gray-50">
  <!-- Header -->
  <header class="bg-white shadow-sm border-b">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <div class="flex items-center space-x-4">
          <Camera class="h-8 w-8 text-primary-600" />
          <h1 class="text-xl font-semibold text-gray-900">Smart City Surveillance</h1>
        </div>

        <div class="flex items-center space-x-4">
          <Button variant="ghost" size="sm" class="relative">
            <Bell class="h-5 w-5" />
            {#if $alerts.alerts.filter(a => a.status === 'pending').length > 0}
              <Badge class="absolute -top-1 -right-1 h-5 w-5 rounded-full p-0 text-xs">
                {$alerts.alerts.filter(a => a.status === 'pending').length}
              </Badge>
            {/if}
          </Button>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" class="relative h-8 w-8 rounded-full">
                <Avatar class="h-8 w-8">
                  <AvatarImage src="" alt={$auth.user?.firstName} />
                  <AvatarFallback>{$auth.user?.firstName?.[0]}{$auth.user?.lastName?.[0]}</AvatarFallback>
                </Avatar>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent class="w-56" align="end" forceMount>
              <DropdownMenuLabel class="font-normal">
                <div class="flex flex-col space-y-1">
                  <p class="text-sm font-medium leading-none">
                    {$auth.user?.firstName} {$auth.user?.lastName}
                  </p>
                  <p class="text-xs leading-none text-muted-foreground">
                    {$auth.user?.email}
                  </p>
                </div>
              </DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem onclick={handleLogout}>
                <LogOut class="mr-2 h-4 w-4" />
                <span>Log out</span>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </div>
  </header>

  <!-- Main Content -->
  <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    {#if isLoading}
      <div class="flex items-center justify-center h-64">
        <div class="animate-spin rounded-full h-32 w-32 border-b-2 border-primary-600"></div>
      </div>
    {:else if error}
      <Card>
        <CardContent class="pt-6">
          <div class="text-center">
            <AlertTriangle class="mx-auto h-12 w-12 text-red-500 mb-4" />
            <h3 class="text-lg font-medium text-gray-900 mb-2">Error Loading Dashboard</h3>
            <p class="text-gray-600 mb-4">{error}</p>
            <Button onclick={loadDashboard}>Retry</Button>
          </div>
        </CardContent>
      </Card>
    {:else}
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Camera Grid -->
        <div class="lg:col-span-2">
          <Card>
            <CardHeader>
              <CardTitle class="flex items-center space-x-2">
                <Video class="h-5 w-5" />
                <span>Live Camera Feeds</span>
              </CardTitle>
              <CardDescription>
                Real-time surveillance feeds from {cameraGrid.length || 0} cameras
              </CardDescription>
            </CardHeader>
            <CardContent>
              {#if cameraGrid.length}
                <div class="grid grid-cols-2 gap-4 aspect-video">
                  {#each cameraGrid as camera}
                      <div class="relative bg-gray-900 rounded-lg overflow-hidden">
                        <div class="absolute inset-0 flex items-center justify-center">
                          <div class="text-center text-white">
                            <Camera class="h-8 w-8 mx-auto mb-2 opacity-50" />
                            <p class="text-sm font-medium">{camera.name}</p>
                            <p class="text-xs opacity-75">{camera.location}</p>
                          </div>
                        </div>
                        <div class="absolute top-2 right-2">
                          <Badge variant={camera.status === 'active' ? 'default' : 'secondary'}>
                            {camera.status}
                          </Badge>
                        </div>
                      </div>
                  {/each}
                </div>
              {:else}
                <div class="text-center py-8">
                  <Camera class="mx-auto h-12 w-12 text-gray-400 mb-4" />
                  <p class="text-gray-500">No cameras available</p>
                </div>
              {/if}
            </CardContent>
          </Card>
        </div>

        <!-- Alerts Panel -->
        <div class="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle class="flex items-center space-x-2">
                <AlertTriangle class="h-5 w-5" />
                <span>Active Alerts</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              {#if $alerts.alerts.length > 0}
                <div class="space-y-3">
                  {#each $alerts.alerts.slice(0, 5) as alert}
                    <div class="p-3 border rounded-lg hover:bg-gray-50 cursor-pointer">
                      <div class="flex items-start justify-between">
                        <div class="flex-1">
                          <h4 class="font-medium text-sm">{alert.title}</h4>
                          <p class="text-xs text-gray-600 mt-1">{alert.description}</p>
                          <div class="flex items-center space-x-2 mt-2">
                            <MapPin class="h-3 w-3 text-gray-400" />
                            <span class="text-xs text-gray-500">{alert.location}</span>
                          </div>
                        </div>
                        <div class="flex flex-col items-end space-y-1">
                          <Badge class={getSeverityColor(alert.severity)}>
                            {alert.severity}
                          </Badge>
                          <Badge class={getStatusColor(alert.status)}>
                            {alert.status}
                          </Badge>
                        </div>
                      </div>
                      <div class="flex items-center mt-2 text-xs text-gray-500">
                        <Clock class="h-3 w-3 mr-1" />
                        {new Date(alert.createdAt).toLocaleTimeString()}
                      </div>
                    </div>
                  {/each}
                </div>
              {:else}
                <div class="text-center py-4">
                  <AlertTriangle class="mx-auto h-8 w-8 text-gray-400 mb-2" />
                  <p class="text-sm text-gray-500">No active alerts</p>
                </div>
              {/if}
            </CardContent>
          </Card>

          <!-- User Info -->
          <Card>
            <CardHeader>
              <CardTitle class="flex items-center space-x-2">
                <Users class="h-5 w-5" />
                <span>User Information</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div class="space-y-3">
                <div class="flex items-center space-x-3">
                  <Avatar>
                    <AvatarImage src="" alt={$auth.user?.firstName} />
                    <AvatarFallback>{$auth.user?.firstName?.[0]}{$auth.user?.lastName?.[0]}</AvatarFallback>
                  </Avatar>
                  <div>
                    <p class="font-medium">{$auth.user?.firstName} {$auth.user?.lastName}</p>
                    <p class="text-sm text-gray-500">{$auth.user?.role}</p>
                  </div>
                </div>
                <div class="text-sm space-y-1">
                  <p><span class="font-medium">Email:</span> {$auth.user?.email}</p>
                  <p><span class="font-medium">Phone:</span> {$auth.user?.phone}</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    {/if}
  </main>
</div> 