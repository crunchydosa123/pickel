import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"
import { Card } from "./ui/card"
import { Settings } from "lucide-react"
import UserCard from "./UserCard"

export function AppSidebar() {
  return (
    <Sidebar>
      <SidebarHeader >
        <Card className="text-left text-2xl px-2 py-3 font-bold rounded-md">Pickel</Card>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup >
          <SidebarGroupLabel>Actions</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>

              <SidebarMenuItem>
                <SidebarMenuButton>Dashboard</SidebarMenuButton>
              </SidebarMenuItem>

              <SidebarMenuItem>
                <SidebarMenuButton>Models</SidebarMenuButton>
              </SidebarMenuItem>

              <SidebarMenuItem>
                <SidebarMenuButton>Deployments</SidebarMenuButton>
              </SidebarMenuItem>

              <SidebarMenuItem>
                <SidebarMenuButton>Ingress</SidebarMenuButton>
              </SidebarMenuItem>

            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        <SidebarGroup />
      </SidebarContent>
      <SidebarFooter >
        <UserCard />
      </SidebarFooter>
    </Sidebar>
  )
}