import { Box } from "@mantine/core";
import { commonLayout } from "@/utils/commonLayout";
import Status from "@/components/index/Status";

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
    </Box>
  )
}

Home.getLayout = commonLayout;