package fr.goralys.app;

import android.os.Build;
import android.os.Bundle;
import android.view.View;

import com.getcapacitor.BridgeActivity;

public class MainActivity extends BridgeActivity {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            getBridge()
                    .getWebView()
                    .setImportantForAutofill(
                            View.IMPORTANT_FOR_AUTOFILL_YES
                    );
        }
    }
}
