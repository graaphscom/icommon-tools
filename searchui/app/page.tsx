import { Icommon } from '@icommon/components/icommon';
import { adn, aws, node } from '@icommon/fontawesome/brands';
import { bell, eye, map } from '@icommon/fontawesome/regular';
import {
  addressBook,
  addressCard,
  alignCenter,
} from '@icommon/fontawesome/solid';
import {
  __3dRotation24px,
  accessibleForward24px,
  accountBalanceSharp24px,
} from '@icommon/material/action';
import {
  addAlert24px,
  notificationImportant24px,
  warningSharp24px,
} from '@icommon/material/alert';
import { moveDownSharp24px } from '@icommon/material/editor';
import { arrowBoth24, blocked24, briefcase24 } from '@icommon/octicons';
import { activityLog, borderDashed, button } from '@icommon/radixui';
import {
  ancientGateFill,
  bankFill,
  homeFill,
} from '@icommon/remixicon/buildings';
import {
  anticlockwiseFill,
  brush2Fill,
  cropFill,
} from '@icommon/remixicon/design';
import { __0Plus, __12Plus, __500px } from '@icommon/unicons/line';
import { angleDoubleDown, apps, bookmark } from '@icommon/unicons/monochrome';
import { redisClient } from '@/db';

export default async function Home() {
  // await (async () => {
  //   const resp = await redisClient.scan();
  //   console.log(resp);
  // })();
  const boxiconsLogos = await import('@icommon/boxicons/logos');
  const boxiconsRegular = await import('@icommon/boxicons/regular');
  const boxiconsSolid = await import('@icommon/boxicons/solid');
  const bytesize = await import('@icommon/bytesize');
  // const fluentui = await import('@icommon/fluentui');
  const fontawesomeSolid = await import('@icommon/fontawesome/solid');
  const fontawesomeRegular = await import('@icommon/fontawesome/regular');
  const fontawesomeBrands = await import('@icommon/fontawesome/brands');
  const octicons = await import('@icommon/octicons');

  return (
    <section>
      <h1>icommon</h1>
      <section>
        <h2>boxicons</h2>
        <section>
          <h3>logos</h3>
          {Object.values(boxiconsLogos).map((l, idx) => (
            <Icommon node={l} key={idx} />
          ))}
        </section>
        <section>
          <h3>regular</h3>
          {Object.values(boxiconsRegular).map((l, idx) => (
            <Icommon node={l} key={idx} />
          ))}
        </section>
        <section>
          <h3>solid</h3>
          {Object.values(boxiconsSolid).map((l, idx) => (
            <Icommon node={l} key={idx} />
          ))}
        </section>
      </section>
      <section>
        <h2>bytesize</h2>
        {Object.values(bytesize).map((l, idx) => (
          <Icommon node={l} key={idx} />
        ))}
      </section>
      <section>
        <h2>fluentui</h2>
        {/*{Object.values(fluentui).map((l, idx) => (*/}
        {/*  <Icommon node={l} key={idx} />*/}
        {/*))}*/}
      </section>
      <section>
        <h2>fontawesome</h2>
        <section>
          <h3>brands</h3>
          {Object.values(fontawesomeBrands).map((l, idx) => (
            <Icommon node={l} key={idx} width="32" height="32" />
          ))}
        </section>
        <section>
          <h3>regular</h3>
          {Object.values(fontawesomeRegular).map((l, idx) => (
            <Icommon node={l} key={idx} width="32" height="32" />
          ))}
        </section>
        <section>
          <h3>solid</h3>
          {Object.values(fontawesomeSolid).map((l, idx) => (
            <Icommon node={l} key={idx} width="32" height="32" />
          ))}
        </section>
      </section>
      <section>
        <h2>material</h2>
      </section>
      <section>
        <h2>octicons</h2>
        {Object.values(octicons).map((l, idx) => (
          <Icommon node={l} key={idx} />
        ))}
      </section>
      <section>
        <h2>radixui</h2>
        <Icommon node={activityLog} />
        <Icommon node={borderDashed} />
        <Icommon node={button} />
      </section>
      <section>
        <h2>remixicon</h2>
        <section>
          <h3>buildings</h3>
          <Icommon node={ancientGateFill} width="32" height="32" />
          <Icommon node={bankFill} width="32" height="32" />
          <Icommon node={homeFill} width="32" height="32" />
        </section>
        <section>
          <h3>design</h3>
          <Icommon node={anticlockwiseFill} width="32" height="32" />
          <Icommon node={brush2Fill} width="32" height="32" />
          <Icommon node={cropFill} width="32" height="32" />
        </section>
      </section>
      <section>
        <h2>unicons</h2>
        <section>
          <h3>line</h3>
          <Icommon node={__0Plus} width="32" height="32" />
          <Icommon node={__12Plus} width="32" height="32" />
          <Icommon node={__500px} width="32" height="32" />
        </section>
        <section>
          <h3>monochrome</h3>
          <Icommon node={angleDoubleDown} width="32" height="32" />
          <Icommon node={apps} width="32" height="32" />
          <Icommon node={bookmark} width="32" height="32" />
        </section>
      </section>
    </section>
  );
}
