import { Box } from "@mantine/core";
import { commonLayout } from "@/utils/commonLayout";
import Status from "@/components/index/Status";
import Capacity from "@/components/index/Capacity";

export async function getStaticProps() {
  return {
    props: {
    }
  }
}

export default function Home() {
  return (
    <Box>
      <Status/>
      <Capacity/>
    </Box>
  )
}

Home.getLayout = commonLayout;