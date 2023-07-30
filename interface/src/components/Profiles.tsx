import { Button, Grid, Typography } from "@mui/material";
import { useVialProfiles } from "../util/hooks";
import { VialProfile, VialProfileCollection } from "../machinepb/machine";
import { KV_KEY_ALL_VIAL_PROFILES } from "../topics_backend/topics_backend";

// I still don't understand the type/value distinction that gives rise to this
// const VialProfileCollectionMethods = VialProfileCollection;

export default function Profiles() {
  const [vialProfiles, setVialProfiles] = useVialProfiles();

  // Experimentation
  // const [vialProfiles, setVialProfiles] = useKeyValueStore<typeof VialProfileCollectionMethods>(
  //   KV_KEY_ALL_VIAL_PROFILES,
  //   VialProfileCollection
  // );
  // const [systemVialProfiles, setSystemVialProfiles] = useSystemVialProfiles();

  // testing
  const setProfiles = () => {
    const msg = vialProfiles ?? VialProfileCollection.create();
    if (!msg.profiles[0]) {
      msg.profiles[0] = VialProfile.create();
    }
    msg.profiles[0].description += "a";
    console.log("setting", msg);
    setVialProfiles(msg);
  };

  return (
    <>
      <Typography variant="h3">vialProfiles</Typography>
      <Button onClick={setProfiles}>Test</Button>
      <Typography variant="body1">{JSON.stringify(vialProfiles, null, 2)}</Typography>
      {/* <Typography variant="h3">systemVialProfiles</Typography> */}
      {/* <Typography variant="body1">{JSON.stringify(systemVialProfiles)}</Typography> */}
    </>
  );
}
