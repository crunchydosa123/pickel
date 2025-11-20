import AddModelCode from '@/components/AddModelCode'
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import React from 'react'
import ModelCode from './code'

type ModelCodeProps = {
  model: any;
};

const ModelTabs = ({ model }: ModelCodeProps) => {
  return (
    <div className="flex justify-between">
      <Tabs defaultValue="account" className="w-full h-screen">
        <TabsList>
          <TabsTrigger value="code">Code</TabsTrigger>
          <TabsTrigger value="password">Deployment</TabsTrigger>
          <TabsTrigger value="observability">Observability</TabsTrigger>
        </TabsList>
        <TabsContent value="code"><ModelCode model={model}/></TabsContent>
        <TabsContent value="password">Change your password here.</TabsContent>
      </Tabs>
      <AddModelCode id={model.id} />
    </div>
  )
}

export default ModelTabs