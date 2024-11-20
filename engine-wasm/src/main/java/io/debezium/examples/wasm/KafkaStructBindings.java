package io.debezium.examples.wasm;

import com.dylibso.chicory.runtime.HostFunction;
import com.dylibso.chicory.runtime.Instance;
import com.dylibso.chicory.wasm.types.ValueType;
import org.apache.kafka.connect.data.Struct;

import java.util.List;

public class KafkaStructBindings {

    public String getString(Struct struct, String fieldName) {
        return struct.getString(fieldName);
    }

    //    public Struct getStruct(Struct struct, String fieldName) {
    //        return struct.getStruct(fieldName);
    //    }

    public static HostFunction getString() {
        return new HostFunction(
                "env",
                "struct_get_string",
                List.of(ValueType.I32, ValueType.I32),
                List.of(ValueType.I32),
                (Instance instance, long... args) -> {
                    System.out.println("Start debugging from here: " + args[0] + " - " + args[1]);
                    return new long[]{ args[0] };
                });
    }
}
