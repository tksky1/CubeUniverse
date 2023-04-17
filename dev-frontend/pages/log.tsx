import { commonLayout } from "@/utils/commonLayout";
import { Box, Tabs, Text } from "@mantine/core";
import { useState, useContext } from "react"
import { DataContext } from "@/components/DataContext";

export default function Log() {
    let data = useContext(DataContext);
    const [activeTab, setActiveTab] = useState<string | null>("core");
    return (
        <Box
            p={"xs"}
        >
            <Tabs value={activeTab} onTabChange={setActiveTab}>
                <Tabs.List>
                    <Tabs.Tab value="core">核心控制器日志</Tabs.Tab>
                    <Tabs.Tab value="back">控制后端日志</Tabs.Tab>
                </Tabs.List>
                <Tabs.Panel
                    value="core"
                    sx={{
                        backgroundColor: "#f7f8fb"
                    }}
                    p={"md"}>
                    <Text>
                        {data && (data as any).Operatorlog ? (data as any).Operatorlog : "暂无日志"}
                    </Text>
                </Tabs.Panel>
                <Tabs.Panel
                    value="back"
                    sx={{
                        backgroundColor: "#f7f8fb"
                    }}
                    p={"md"}>
                    <Text>
                        {data && (data as any).Backendlog ? (data as any).Backendlog : "暂无日志"}
                    </Text>
                </Tabs.Panel>
            </Tabs>
        </Box>
    )
}

Log.getLayout = commonLayout;