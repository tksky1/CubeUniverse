import MyTable from "@/components/MyTable";
import { commonLayout } from "@/utils/commonLayout"
import { Box } from "@mantine/core";

export default function Page1() {
    return (
        <Box
            p={"md"}
            sx={theme => ({
            })}>
            <MyTable/>
        </Box>
    )
}
Page1.getLayout = commonLayout;